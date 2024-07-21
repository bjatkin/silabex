from nltk.corpus import cmudict
from g2p_en import G2p
import polars as pl
from datetime import datetime
import sys
import re

def valid_sylables(sylable):
    for sylable in sylables:
        vowels = [ phoneme for phoneme in sylable if re.match(".*\d+", phoneme) ]
        if len(vowels) != 1:
            return False

    return True

def valid_label(phonemes, label):
    # check for the valid controll commands
    if label in ["mono", "skip", "exit", "quit"]:
        return True

    simple_phonemes = (
        "".join(phonemes)
        .replace("0", "")
        .replace("1", "")
        .replace("2", "")
        .lower()
    )
    simple_label = label.replace("/", "")

    return simple_phonemes == simple_label
    

def build_sylables(phonemes, label):
    if label == "mono":
        return [phonemes]
    
    label_offset = 0
    sylables = []
    sylable = []
    for phoneme in phonemes:
        sylable.append(phoneme)
        
        label_offset += len(
            phoneme
            .removesuffix("0")
            .removesuffix("1")
            .removesuffix("2")
        )
        
        if label_offset >= len(label) or label[label_offset] == "/":
            label_offset += 1
            sylables.append(sylable)
            sylable = []


    return sylables


with open("data/dictionary/popular.txt", "r") as f:
    popular_words = pl.DataFrame(f.read().split("\n"), schema={"word": pl.String})
    
# get config data (TODO: get this from an env file or something?)
batch_size = sys.argv[1] if len(sys.argv) >= 2 else 25
label_path = "data/sylables/sylables.csv"
checkpoint_path = f"data/sylables/sylables_checkpoint_{datetime.utcnow()}.csv"
save_labels = True

g2p = G2p()
labeled_data = pl.read_csv(
    label_path,
    schema={"word": pl.String, "phonemes": pl.String, "sylables": pl.String, "labeling_time": pl.Float64},
)


# save a checkpoint before labeling more data
if labeled_data.shape[0] > 0 and save_labels:
    labeled_data.write_csv(checkpoint_path)
    print(f"info: writing checkpoint file [{checkpoint_path}]")


print(f"loading a batch of {batch_size} words for labeling")
for row in popular_words.sample(batch_size).iter_rows(named=True):
    word = row["word"]
    phonemes = g2p(word)

    start_labeling = datetime.now()
    print("Data:", word, "|", phonemes)
    is_valid = False
    while not is_valid:
        print("> ", end="")
        
        label = input()
        is_valid = valid_label(phonemes, label)
        
        if not is_valid:
            print("warning: invalid label does not match word phonemes, try again")
            continue

        if label in ["skip", "exit", "quit"]:
            break
            
        sylables = build_sylables(phonemes, label)
        is_valid = valid_sylables(sylables)

        if not is_valid:
            print("warning: invalid label, sylable has the wrong number of vowels, try again")
            continue

    if label == "skip":
        print(f"skipping {word}")
        continue

    if label in ["exit", "quit"]:
        if save_labels:
            print(f"info: saving labels and then exiting [{label_path}]")
            labeled_data.write_csv(label_path)
        break
    
        
    elapsed_seconds = (datetime.now()-start_labeling).total_seconds()
    print(f"info: labeled in {elapsed_seconds:.2f} seconds")

    row = pl.DataFrame(
        {
            "word": word,
            "phonemes": str(phonemes),
            "sylables": str(sylables),
            "labeling_time": elapsed_seconds,
        },
    )
    labeled_data = pl.concat([labeled_data, row])

    if save_labels:
        labeled_data.write_csv(label_path)
        print(f"info: labels saved [{label_path}]")
