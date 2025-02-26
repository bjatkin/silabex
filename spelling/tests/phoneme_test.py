import phoneme as ph
import pytest


@pytest.mark.parametrize(
    "text,want",
    [
        ("this is a test", ["this", "is", "a", "test"]),
        (
            "this is (another, test) to run $100",
            ["this", "is", "another", "test", "to", "run"],
        ),
    ],
)
def test_split_words(text: str, want: list[str]):
    """
    Test split_words to make sure larger texts are correctly split into smaller words
    """

    got = ph.split_words(text)
    assert got == want


def build_word(word: str, arpa_word: list[list[str]]) -> ph.Word:
    syllables: list[ph.Syllable] = []
    for syllable in arpa_word:
        phonemes = [ph.Phoneme.from_string(p) for p in syllable]
        arpa_syllable = ph.Syllable(phonemes=phonemes)
        syllables.append(arpa_syllable)

    return ph.Word(word=word, arpabet_word=syllables)


@pytest.mark.parametrize(
    "word, want",
    [
        # long words
        (
            "championship",
            build_word(
                "championship",
                [["CH", "AE1", "M"], ["P", "IY0"], ["AH0", "N"], ["SH", "IH2", "P"]],
            ),
        ),
        # words with apostraphies
        (
            "tournament's",
            ph.Word(
                word="tournament's",
                arpabet_word=[
                    ph.Syllable(
                        phonemes=[
                            ph.Phoneme.from_string("T"),
                            ph.Phoneme.from_string("UH1"),
                            ph.Phoneme.from_string("R"),
                        ]
                    ),
                    ph.Syllable(
                        phonemes=[
                            ph.Phoneme.from_string("N"),
                            ph.Phoneme.from_string("AH0"),
                        ]
                    ),
                    ph.Syllable(
                        phonemes=[
                            ph.Phoneme.from_string("M"),
                            ph.Phoneme.from_string("AH0"),
                            ph.Phoneme.from_string("N"),
                            ph.Phoneme.from_string("T"),
                            ph.Phoneme(
                                phoneme=ph.ARPAbet.S, sonority=-1, is_vowel=False
                            ),
                        ]
                    ),
                ],
            ),
        ),
        ("can't", build_word("can't", [["K", "AE1", "N", "T"]])),
        ("won't", build_word("won't", [["W", "OW1", "N", "T"]])),
        ("shouldn't", build_word("shouldn't", [["SH", "UH1"], ["D", "AH0", "N", "T"]])),
        # words with inner capital letters
        (
            "McDonald",
            build_word(
                "McDonald", [["M", "IH0", "K"], ["D", "AA1"], ["N", "AH0", "L", "D"]]
            ),
        ),
        (
            "D'Angelo",
            build_word(
                "D'Angelo", [["D", "IY1"], ["AE1", "N"], ["JH", "AH0"], ["L", "OW2"]]
            ),
        ),
        ("iPhone", build_word("iPhone", [["AY1"], ["F", "OW1", "N"]])),
        ("eBay", build_word("eBay", [["IY1"], ["B", "EY1"]])),
        (
            "CoPilot",
            build_word("CoPilot", [["K", "OW1"], ["P", "AY1"], ["L", "AH0", "T"]]),
        ),
        # words with adjacent vowels
        ("fire", build_word("fire", [["F", "AY1"], ["ER0"]])),
        ("hour", build_word("hour", [["AW1"], ["ER0"]])),
        ("player", build_word("player", [["P", "L", "EY1"], ["ER0"]])),
        # words with silent letters
        ("knight", build_word("knight", [["N", "AY1", "T"]])),
        ("aisle", build_word("aisle", [["AY1", "L"]])),
        ("subtle", build_word("subtle", [["S", "AH1"], ["T", "AH0", "L"]])),
        # words with silent letters
        ("knight", build_word("knight", [["N", "AY1", "T"]])),
        ("aisle", build_word("aisle", [["AY1", "L"]])),
        ("subtle", build_word("subtle", [["S", "AH1"], ["T", "AH0", "L"]])),
        # words with schwa
        (
            "separate",
            build_word("separate", [["S", "EH1"], ["P", "ER0"], ["EY2", "T"]]),
        ),
        ("camera", build_word("camera", [["K", "AE1"], ["M", "ER0"], ["AH0"]])),
        ("family", build_word("family", [["F", "AE1"], ["M", "AH0"], ["L", "IY0"]])),
        # words with consonant clusters
        (
            "strength",
            ph.Word(
                word="strength",
                arpabet_word=[
                    ph.Syllable(
                        phonemes=[
                            ph.Phoneme(
                                phoneme=ph.ARPAbet.S, sonority=-1, is_vowel=False
                            ),
                            ph.Phoneme.from_string("T"),
                            ph.Phoneme.from_string("R"),
                            ph.Phoneme.from_string("EH1"),
                            ph.Phoneme.from_string("NG"),
                            ph.Phoneme.from_string("K"),
                            ph.Phoneme.from_string("TH"),
                        ]
                    )
                ],
            ),
        ),
        (
            "sprints",
            ph.Word(
                word="sprints",
                arpabet_word=[
                    ph.Syllable(
                        phonemes=[
                            ph.Phoneme(
                                phoneme=ph.ARPAbet.S, sonority=-1, is_vowel=False
                            ),
                            ph.Phoneme.from_string("P"),
                            ph.Phoneme.from_string("R"),
                            ph.Phoneme.from_string("IH1"),
                            ph.Phoneme.from_string("N"),
                            ph.Phoneme.from_string("T"),
                            ph.Phoneme(
                                phoneme=ph.ARPAbet.S, sonority=-1, is_vowel=False
                            ),
                        ]
                    )
                ],
            ),
        ),
        ("through", build_word("through", [["TH", "R", "UW1"]])),
        (
            "blasts",
            ph.Word(
                word="blasts",
                arpabet_word=[
                    ph.Syllable(
                        phonemes=[
                            ph.Phoneme.from_string("B"),
                            ph.Phoneme.from_string("L"),
                            ph.Phoneme.from_string("AE1"),
                            ph.Phoneme.from_string("S"),
                            ph.Phoneme.from_string("T"),
                            ph.Phoneme(
                                phoneme=ph.ARPAbet.S, sonority=-1, is_vowel=False
                            ),
                        ]
                    )
                ],
            ),
        ),
        # compound words
        (
            "wallpaper",
            build_word("wallpaper", [["W", "AO1", "L"], ["P", "EY2"], ["P", "ER0"]]),
        ),
        ("football", build_word("football", [["F", "UH1", "T"], ["B", "AO2", "L"]])),
        (
            "bookshelf",
            build_word("bookshelf", [["B", "UH1", "K"], ["SH", "EH2", "L", "F"]]),
        ),
        # unusual stress patterns
        ("guitar", build_word("guitar", [["G", "IH0"], ["T", "AA1", "R"]])),
        ("police", build_word("police", [["P", "AH0"], ["L", "IY1", "S"]])),
        # loan words
        (
            "croissant",
            build_word("croissant", [["K", "W", "AA2"], ["S", "AA1", "N", "T"]]),
        ),
        ("cliche", build_word("cliche", [["K", "L", "IY0"], ["SH", "EY1"]])),
        (
            "karaoke",
            build_word("karaoke", [["K", "EH2"], ["R", "IY0"], ["OW1"], ["K", "IY0"]]),
        ),
        # exceptions to the sonority sequencing principal
        ("psalm", build_word("psalm", [["S", "AA1", "L", "M"]])),
        ("sphere", build_word("sphere", [["S", "F", "IH1", "R"]])),
        (
            "split",
            ph.Word(
                word="split",
                arpabet_word=[
                    ph.Syllable(
                        phonemes=[
                            ph.Phoneme(
                                phoneme=ph.ARPAbet.S, sonority=-1, is_vowel=False
                            ),
                            ph.Phoneme.from_string("P"),
                            ph.Phoneme.from_string("L"),
                            ph.Phoneme.from_string("IH1"),
                            ph.Phoneme.from_string("T"),
                        ]
                    )
                ],
            ),
        ),
        (
            "sprint",
            ph.Word(
                word="sprint",
                arpabet_word=[
                    ph.Syllable(
                        phonemes=[
                            ph.Phoneme(
                                phoneme=ph.ARPAbet.S, sonority=-1, is_vowel=False
                            ),
                            ph.Phoneme.from_string("P"),
                            ph.Phoneme.from_string("R"),
                            ph.Phoneme.from_string("IH1"),
                            ph.Phoneme.from_string("N"),
                            ph.Phoneme.from_string("T"),
                        ]
                    )
                ],
            ),
        ),
        (
            "sphinx",
            ph.Word(
                word="sphinx",
                arpabet_word=[
                    ph.Syllable(
                        phonemes=[
                            ph.Phoneme.from_string("S"),
                            ph.Phoneme.from_string("F"),
                            ph.Phoneme.from_string("IH1"),
                            ph.Phoneme.from_string("NG"),
                            ph.Phoneme.from_string("K"),
                            ph.Phoneme(
                                phoneme=ph.ARPAbet.S, sonority=-1, is_vowel=False
                            ),
                        ]
                    )
                ],
            ),
        ),
        ("thrive", build_word("thrive", [["TH", "R", "AY1", "V"]])),
        ("free", build_word("free", [["F", "R", "IY1"]])),
        ("twelfth", build_word("twelfth", [["T", "W", "EH1", "L", "F", "TH"]])),
        ("lengths", build_word("lengths", [["L", "EH1", "NG", "K", "TH", "S"]])),
    ],
)
def test_word_init(word: str, want: ph.Word):
    got = ph.Word.from_string(word)
    assert got == want
