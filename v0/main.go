package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/bjatkins/silabex/font"
)

type Fragment interface {
	Render(string) string
}

type Path struct {
	// TODO: should I break this down further?
	d string
}

func (p Path) Render(style string) string {
	return fmt.Sprintf("<path style=\"%s\" d=\"%s\" />", style, p.d)
}

type Circle struct {
	cx float64
	cy float64
	r  float64
}

func (c Circle) Render(style string) string {
	return fmt.Sprintf("<circle style=\"%s\" cx=\"%f\" cy=\"%f\" r=\"%f\" />", style, c.cx, c.cy, c.r)
}

var (
	LeftVowel       Fragment = Path{d: "M 6.7052574,168.07537 V 1.6923473"}
	RightVowel      Fragment = Path{d: "M 163.79351,168.04203 V 1.6391129"}
	TopVowel        Fragment = Path{d: "M 2.5513018,5.6210201 H 167.66428"}
	BottomVowel     Fragment = Path{d: "M 2.5222205,164.2161 H 167.84669"}
	CircVowel       Fragment = Circle{cx: 84.353699, cy: 85.079872, r: 29.004452}
	Pre5            Fragment = Path{d: "M 21.385177,145.54081 H 77.238139"}
	Pre4            Fragment = Path{d: "M 30.53434,130.0988 H 71.368991"}
	Pre3            Fragment = Path{d: "M 30.241845,85.668013 71.147841,84.642862"}
	Pre2            Fragment = Path{d: "M 30.453647,41.942372 H 71.445106"}
	Pre1            Fragment = Path{d: "M 23.398572,24.177821 H 77.192716"}
	PreSlash        Fragment = Path{d: "m 38.270476,74.601842 11.259272,19.50163 M 54.042275,74.566657 65.301547,94.06829"}
	PreRight        Fragment = Path{d: "M 67.518632,134.19146 V 38.029759"}
	PreLeft         Fragment = Path{d: "M 34.471006,134.44306 V 37.737942"}
	Post5           Fragment = Path{d: "M 92.696549,145.9062 H 149.17321"}
	Post4           Fragment = Path{d: "M 97.649021,130.78911 H 139.58542"}
	Post3           Fragment = Path{d: "M 98.016258,84.810774 H 139.8166"}
	Post2           Fragment = Path{d: "M 98.599368,42.587475 H 140.08131"}
	Post1           Fragment = Path{d: "M 91.801656,23.957394 H 150.67206"}
	PostCirc        Fragment = Circle{cx: 118.22502, cy: 84.42823, r: 16.378378}
	PostMidSlash    Fragment = Path{d: "m 106.22476,73.842705 11.25927,19.50163 m 4.51253,-19.536815 11.25927,19.501633"}
	PostTopSlash    Fragment = Path{d: "m 105.5183,13.803039 11.25927,19.501629 m 4.51253,-19.536814 11.25927,19.501634"}
	PostBottomSlash Fragment = Path{d: "m 107.11249,136.99923 11.25927,19.50163 m 4.51253,-19.53682 11.25927,19.50163"}
	PostLeft        Fragment = Path{d: "M 101.88452,134.36309 V 38.817728"}
	PostRight       Fragment = Path{d: "M 136.2636,134.87325 V 39.085363"}
	Full5           Fragment = Path{d: "M 37.442724,145.58561 H 131.98925"}
	Full4           Fragment = Path{d: "M 38.968265,129.6315 H 131.40105"}
	Full3           Fragment = Path{d: "M 39.611371,84.242222 H 130.62497"}
	Full2           Fragment = Path{d: "M 39.095843,41.865901 131.28094,41.744027"}
	Full1           Fragment = Path{d: "M 38.331702,23.827328 132.77855,23.636803"}
	FullLeft        Fragment = Path{d: "M 42.901158,133.512 V 38.442161"}
	FullRight       Fragment = Path{d: "M 127.56835,132.7553 V 37.524682"}
)

type Cluster []Fragment

var (
	ACluster  = Cluster{LeftVowel}
	ALCluster = Cluster{LeftVowel, RightVowel}

	ECluster = Cluster{TopVowel}

	ICluster   = Cluster{TopVowel, RightVowel}
	IyeCluster = Cluster{TopVowel, LeftVowel, RightVowel, BottomVowel}

	OCluster  = Cluster{BottomVowel}
	OOCluster = Cluster{LeftVowel, BottomVowel}

	UCluster = Cluster{RightVowel}

	PostNCluster = Cluster{Post2, PostLeft, PostRight}

	PreTCluster = Cluster{Pre2, PreRight}
	PreMCluster = Cluster{PreRight, Pre4, Pre3}
	PreNCluster = Cluster{Pre2, PreRight, Pre3, Pre4}

	FullTCluster = Cluster{Full2, FullRight}
	FullNCluster = Cluster{Full2, FullRight, Full3, Full4}
)

type Rune []Cluster

func (i Rune) SVG(size int, style string) string {
	fragments := ""
	for _, cluster := range i {
		for _, fragment := range cluster {
			fragments += fragment.Render(style)
		}
	}

	return "<svg " +
		"version=\"1.1\" " +
		fmt.Sprintf("width=\"%d\" ", size) +
		fmt.Sprintf("height=\"%d\" ", size) +
		"viewBox=\"0 0 170 170\" " +
		"xmlns=\"http://www.w3.org/2000/svg\">" +
		fragments +
		"</svg>"
}

func main2() {
	tiny := Rune{PreTCluster, IyeCluster, PostNCluster}
	ti := Rune{FullTCluster, IyeCluster}
	ny := Rune{FullNCluster, ICluster}
	moon := Rune{PreMCluster, OOCluster, PostNCluster}

	style := "fill:none;stroke:#000000;stroke-width:10;stroke-dasharray:none;stroke-opacity:1"
	size := 24

	tinySVG := tiny.SVG(size, style)
	tiSVG := ti.SVG(size, style)
	nySVG := ny.SVG(size, style)
	moonSVG := moon.SVG(size, style)

	os.WriteFile("test.html", []byte("<!DOCTYPE html>"+
		"<html>"+
		"<head>"+
		"<title>Silabex</title>"+
		"<style>"+
		"svg { padding-left: 2px; }"+
		"</style>"+
		"</head>"+
		"<body>"+
		"<div>"+tinySVG+"</div>"+
		"<div>"+tiSVG+"</div>"+
		"<div>"+nySVG+"</div>"+
		"<div>"+moonSVG+"</div>"+
		"</body>"+
		"</html>"), 0o0655)
}

func main() {
	font, err := font.New("examples/basic_font_2.svg", "font/derived")
	if err != nil {
		fmt.Println("font err:", err)
		return
	}

	chars, err := font.NewChars("KOE/KWOE/KROE/KHOE/ROE/HOE")
	if err != nil {
		fmt.Println("char err:", err)
		return
	}

	for _, char := range chars {
		filepath := fmt.Sprintf("gen/char_%s.svg", char.Name())
		os.WriteFile(filepath, []byte(char.SVG()), 0o0655)
		err := openInBrowser(filepath)
		if err != nil {
			fmt.Println("browser err:", err)
			return
		}
	}
}

func openInBrowser(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", path)
	case "darwin":
		cmd = exec.Command("open", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	return cmd.Run()
}
