from pydantic import BaseModel
from enum import Enum
import pronouncing


def split_words(text: str) -> list[str]:
    """
    Split text into word-ish segments based on normal word seperators.

    Args:
        text (str): The input text to parse into words.

    Returns:
        list[str]: A list of words.
    """

    valid = "abcdefghijklmnopqrstuvwxyz'"
    valid += valid.upper()
    boundry = "()/,:;.?!*$ "

    all_words = []
    current_word = ""

    for char in text:
        if char in valid:
            current_word += char
            continue

        if char in boundry and current_word != "":
            all_words.append(current_word)

        current_word = ""

    if current_word != "":
        all_words.append(current_word)

    return all_words


class ARPAbet(Enum):
    # Vowels
    AA = "AA"  # odd
    AE = "AE"  # at
    AH = "AH"  # hut
    AO = "AO"  # ought
    AW = "AW"  # cow
    AY = "AY"  # hide
    EH = "EH"  # end
    ER = "ER"  # hurt
    EY = "EY"  # ate
    IH = "IH"  # it
    IY = "IY"  # eat
    OW = "OW"  # oat
    OY = "OY"  # toy
    UH = "UH"  # book
    UW = "UW"  # two

    # Consonants
    B = "B"  # bat
    CH = "CH"  # chin
    D = "D"  # day
    DH = "DH"  # they
    F = "F"  # fat
    G = "G"  # go
    HH = "HH"  # hat
    JH = "JH"  # judge
    K = "K"  # kite
    L = "L"  # let
    M = "M"  # man
    N = "N"  # nap
    NG = "NG"  # sing
    P = "P"  # pet
    R = "R"  # red
    S = "S"  # sit
    SH = "SH"  # she
    T = "T"  # top
    TH = "TH"  # thin
    V = "V"  # van
    W = "W"  # win
    Y = "Y"  # yes
    Z = "Z"  # zoo
    ZH = "ZH"  # measure


class Phoneme(BaseModel):
    """
    An ARPAbet phneme
    """

    phoneme: ARPAbet | str
    sonority: int
    is_vowel: bool

    @classmethod
    def from_string(cls, phoneme: str) -> "Phoneme":
        if phoneme == "<START>" or phoneme == "<END>":
            return cls(phoneme=phoneme, sonority=-1, is_vowel=False)

        phoneme = phoneme.strip("012").upper()
        arpa_phoneme = ARPAbet[phoneme]

        sonority_map = {
            # Vowels (highest sonority)
            ARPAbet.AA: 6,
            ARPAbet.AE: 6,
            ARPAbet.AH: 6,
            ARPAbet.AO: 6,
            ARPAbet.AW: 6,
            ARPAbet.AY: 6,
            ARPAbet.EH: 6,
            ARPAbet.ER: 6,
            ARPAbet.EY: 6,
            ARPAbet.IH: 6,
            ARPAbet.IY: 6,
            ARPAbet.OW: 6,
            ARPAbet.OY: 6,
            ARPAbet.UH: 6,
            ARPAbet.UW: 6,
            # Glides (semi-vowels)
            ARPAbet.W: 5,
            ARPAbet.Y: 5,
            # Liquids
            ARPAbet.L: 4,
            ARPAbet.R: 4,
            # Nasals
            ARPAbet.M: 3,
            ARPAbet.N: 3,
            ARPAbet.NG: 3,
            # Fricatives
            ARPAbet.DH: 2,
            ARPAbet.F: 2,
            ARPAbet.V: 2,
            ARPAbet.TH: 2,
            ARPAbet.S: 2,
            ARPAbet.Z: 2,
            ARPAbet.SH: 2,
            ARPAbet.ZH: 2,
            ARPAbet.HH: 2,
            # Plosives (Stops)
            ARPAbet.B: 1,
            ARPAbet.D: 1,
            ARPAbet.G: 1,
            ARPAbet.K: 1,
            ARPAbet.P: 1,
            ARPAbet.T: 1,
            ARPAbet.CH: 1,
            ARPAbet.JH: 1,
        }

        sonority = sonority_map[arpa_phoneme]

        vowels = [
            ARPAbet.AA,
            ARPAbet.AE,
            ARPAbet.AH,
            ARPAbet.AO,
            ARPAbet.AW,
            ARPAbet.AY,
            ARPAbet.EH,
            ARPAbet.ER,
            ARPAbet.EY,
            ARPAbet.IH,
            ARPAbet.IY,
            ARPAbet.OW,
            ARPAbet.OY,
            ARPAbet.UH,
            ARPAbet.UW,
        ]

        is_vowel = arpa_phoneme in vowels

        return cls(phoneme=arpa_phoneme, sonority=sonority, is_vowel=is_vowel)

    def __repr__(self):
        if isinstance(self.phoneme, str):
            return self.phoneme

        if self.sonority == -1:
            return "-" + self.phoneme.value

        return self.phoneme.value


StartPhoneme = Phoneme.from_string("<START>")
EndPhoneme = Phoneme.from_string("<END>")


class Syllable(BaseModel):
    phonemes: list[Phoneme] = []

    def insert_phoneme(self, phoneme: Phoneme):
        self.phonemes.insert(0, phoneme)

    def contains_vowel(self):
        for p in self.phonemes:
            if p.is_vowel:
                return True

        return False

    def is_empty(self) -> bool:
        return len(self.phonemes) == 0

    def __repr__(self) -> str:
        return " ".join([repr(p) for p in self.phonemes])


class Word(BaseModel):
    word: str
    arpabet_word: list[Syllable]

    @classmethod
    def from_string(cls, word: str) -> "Word":
        # check if there's an inner capital like in words McDonald or iPhone
        # this check dosen't catch single letter words but that's fine since
        # those words will just be a single syllable anyway
        if word[1:].lower() != word[1:]:
            chars = [" " + char if char.isupper() else char for char in word]
            sub_words = "".join(chars).split()

            arpabet_word: list[Syllable] = []
            for sub_word in sub_words:
                # this is important for words like D'Angelo
                sub_word = sub_word.removesuffix("'")
                arpabet_word.extend(get_syllables(sub_word))

            return cls(word=word, arpabet_word=arpabet_word)

        arpabet_word = get_syllables(word)
        return cls(word=word, arpabet_word=arpabet_word)

    def __repr__(self):
        arpa_str = " / ".join([repr(s) for s in self.arpabet_word])
        return f"{self.word} ({arpa_str})"


def get_phonemes(word: str) -> list[Phoneme]:
    phonemes = pronouncing.phones_for_word(word)
    if phonemes is None or len(phonemes) == 0:
        raise ValueError(f"unknown word '{word}'")

    phonemes = phonemes[0].split()
    phonemes = (
        [StartPhoneme] + [Phoneme.from_string(p) for p in phonemes] + [EndPhoneme]
    )

    # pre-process sonorities to handle special case for 'S'
    phonemen_tripplet = zip(phonemes[:-2], phonemes[1:-1], phonemes[2:])
    for left, center, right in phonemen_tripplet:
        # special case for english cause it's weird
        if (
            center.phoneme == ARPAbet.S
            and left.sonority < center.sonority
            and center.sonority > right.sonority
        ):
            center.sonority = -1

    return phonemes


def get_syllables(word: str) -> list[Syllable]:
    phonemes = get_phonemes(word)

    syllable = Syllable()
    syllables: list[Syllable] = []
    phonemen_tripplet = zip(phonemes[:-2], phonemes[1:-1], phonemes[2:])
    for left, center, right in reversed(list(phonemen_tripplet)):
        if center.is_vowel and syllable.contains_vowel():
            syllables.insert(0, syllable)
            syllable = Syllable()
            syllable.insert_phoneme(center)
            continue

        syllable.insert_phoneme(center)
        if (
            left.sonority > center.sonority
            and center.sonority < right.sonority
            and syllable.contains_vowel()
        ):
            syllables.insert(0, syllable)
            syllable = Syllable()

    if not syllable.is_empty():
        syllables.insert(0, syllable)

    return syllables
