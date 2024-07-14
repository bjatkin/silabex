package font

import (
	"fmt"
	"os"
	"strings"
)

type builderError struct {
	fileName  string
	errors    []*lineError
	maxErrors int
}

func newBuilderError(fileName string, errors []*lineError) *builderError {
	return &builderError{
		fileName:  fileName,
		errors:    errors,
		maxErrors: 10,
	}
}

func (b *builderError) Error() string {
	dir, _ := os.Getwd()
	fileName := strings.Replace(b.fileName, dir, ".", 1)
	errors := []string{}
	for i, e := range b.errors {
		if i > b.maxErrors {
			errors = append(errors, fmt.Sprintf("max errors reached, there were %d more errors", len(b.errors)-i))
			break
		}
		errors = append(errors, fmt.Sprintf("%s:%d %s", fileName, e.lineNumber, e.message))
	}

	return strings.Join(errors, "\n")
}

type lineError struct {
	lineNumber int
	message    string
}

type builderTransform int

const (
	moveDown builderTransform = iota
	moveUp
	moveLeft
	moveRight
	mirriorX
	mirriorY
)

type builderExpr struct {
	name       string
	transforms []builderTransform
}

func (b *builderExpr) build(clusters map[string]*cluster) *cluster {
	cluster := clusters[b.name].clone()

	for _, transform := range b.transforms {
		switch transform {
		case moveDown:
			switch cluster.slot {
			case vowelSlot:
				cluster.translate(0, 880)
			case initialSlot:
				cluster.translate(0, 140)
			}
		case moveUp:
			switch cluster.slot {
			case vowelSlot:
				cluster.translate(0, -880)
			case initialSlot:
				cluster.translate(0, -140)
			}
		case moveLeft:
			if cluster.slot == vowelSlot {
				cluster.translate(-880, 0)
			}
		case moveRight:
			if cluster.slot == vowelSlot {
				cluster.translate(880, 0)
			}
		case mirriorX:
			cluster.mirriorX()
		case mirriorY:
			cluster.mirriorY()
		}
	}

	return cluster
}

type builderStmt struct {
	name  string
	exprs []builderExpr
}

func (b *builderStmt) build(clusters map[string]*cluster) *cluster {
	clusterList := []*cluster{}
	for _, expr := range b.exprs {
		clusterList = append(clusterList, expr.build(clusters))
	}

	return mergeClusters(clusterList...)
}

type clusterBuilder struct {
	// metadata *metadata
	initial []builderStmt
	vowel   []builderStmt
	final   []builderStmt
}

type sectionType int

const (
	sectionNone sectionType = iota
	sectionUnknown
	sectionVowel
	sectionInitial
	sectionFinal
)

func newClusterBuilder(file string) (*clusterBuilder, error) {
	raw, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	errors := []*lineError{}
	builder := &clusterBuilder{}

	section := sectionNone
	for i, line := range strings.Split(string(raw), "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		switch strings.TrimSpace(line) {
		case "VOWEL":
			if section != sectionNone {
				errors = append(errors, &lineError{
					lineNumber: i + 1,
					message:    "found the VOWEL section but it was not the first section in the file",
				})
			}

			section = sectionVowel
		case "INITIAL":
			if section != sectionVowel {
				errors = append(errors, &lineError{
					lineNumber: i + 1,
					message:    "found the INITIAL section but it was not preceeded by the VOWEL section",
				})
			}

			section = sectionInitial
		case "FINAL":
			if section != sectionInitial {
				errors = append(errors, &lineError{
					lineNumber: i + 1,
					message:    "found the FINAL section but it was not preceeded by the INITIAL section",
				})
			}

			section = sectionFinal
		default:
			s, err := parseLine(section, i+1, line)
			if err != nil {
				errors = append(errors, err)
				continue
			}

			err = builder.addStmt(section, i+1, *s)
			if err != nil {
				errors = append(errors, err)
			}
		}
	}

	if len(errors) > 0 {
		return nil, newBuilderError(file, errors)
	}

	return builder, nil
}

func (b *clusterBuilder) addStmt(section sectionType, lineNumber int, stmt builderStmt) *lineError {
	switch section {
	case sectionNone:
		return &lineError{
			lineNumber: lineNumber,
			message:    "no section was set, a section must be specified before providing cluster derivations",
		}
	case sectionUnknown:
		return &lineError{
			lineNumber: lineNumber,
			message:    "encountered an unknown cluster section",
		}
	case sectionVowel:
		b.vowel = append(b.vowel, stmt)
	case sectionInitial:
		b.initial = append(b.initial, stmt)
	case sectionFinal:
		b.final = append(b.final, stmt)
	}

	return nil
}

func parseLine(section sectionType, lineNumber int, line string) (*builderStmt, *lineError) {
	parts := strings.Split(line, "|")
	if len(parts) != 2 {
		return nil, &lineError{
			lineNumber: lineNumber,
			message:    "missing | seperator for line",
		}
	}

	name := strings.TrimSpace(parts[0])
	switch section {
	case sectionInitial:
		err := validateInitialName(lineNumber, name)
		if err != nil {
			return nil, err
		}
	case sectionVowel:
		err := validateVowelName(lineNumber, name)
		if err != nil {
			return nil, err
		}
	case sectionFinal:
		err := validateFinalName(lineNumber, name)
		if err != nil {
			return nil, err
		}
	default:
		return nil, &lineError{
			lineNumber: lineNumber,
			message:    "can not valid line because the section is unknown",
		}
	}

	exprs := []builderExpr{}
	for _, cluster := range strings.Split(parts[1], " ") {
		if cluster == "" {
			continue
		}

		parts := strings.Split(cluster, ".")

		switch section {
		case sectionInitial:
			err := validateInitialName(lineNumber, parts[0])
			if err != nil {
				return nil, err
			}
		case sectionVowel:
			err := validateVowelName(lineNumber, parts[0])
			if err != nil {
				return nil, err
			}
		case sectionFinal:
			err := validateFinalName(lineNumber, parts[0])
			if err != nil {
				return nil, err
			}
		default:
			return nil, &lineError{
				lineNumber: lineNumber,
				message:    "can not valid line because the section is unknown",
			}
		}

		if len(parts) == 1 {
			exprs = append(exprs, builderExpr{name: parts[0]})
			continue
		}

		cmds := []builderTransform{}
		for _, cmd := range parts[1] {
			switch cmd {
			case 'd':
				cmds = append(cmds, moveDown)
			case 'u':
				cmds = append(cmds, moveUp)
			case 'l':
				cmds = append(cmds, moveLeft)
			case 'r':
				cmds = append(cmds, moveRight)
			case 'x':
				cmds = append(cmds, mirriorX)
			case 'y':
				cmds = append(cmds, mirriorY)
			default:
				return nil, &lineError{
					lineNumber: lineNumber,
					message:    fmt.Sprintf("unknown command '%s' for cluster '%s'", string(cmd), parts[0]),
				}
			}
		}

		exprs = append(exprs, builderExpr{
			name:       parts[0],
			transforms: cmds,
		})
	}

	return &builderStmt{
		name:  name,
		exprs: exprs,
	}, nil
}

func validateVowelName(lineNumber int, name string) *lineError {
	return validateName(lineNumber, name, "vowel", "AOEU")
}

func validateInitialName(lineNumber int, name string) *lineError {
	return validateName(lineNumber, name, "initial", "KTWPRH")
}

func validateFinalName(lineNumber int, name string) *lineError {
	return validateName(lineNumber, name, "final", "BPGLST")
}

func validateName(lineNumber int, name, checkName, check string) *lineError {
	order := 0
	for i, r := range name {
		hilight := fmt.Sprintf("%s[%s]", name[:i], name[i:i+1])
		if i+1 < len(name) {
			hilight = fmt.Sprintf("%s[%s]%s", name[:i], name[i:i+1], name[i+1:])

		}

		if !strings.Contains(check, string(r)) {
			return &lineError{
				lineNumber: lineNumber,
				message:    fmt.Sprintf("invalid %s name contains an unknown character '%s'", checkName, hilight),
			}
		}

		index := strings.Index(check, string(r))
		if index < order {
			return &lineError{
				lineNumber: lineNumber,
				message:    fmt.Sprintf("invalid %s order '%s' vowels must follow the order %s", checkName, hilight, check),
			}
		}

		order = index
	}

	return nil
}
