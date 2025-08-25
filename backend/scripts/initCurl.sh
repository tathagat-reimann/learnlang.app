#!/usr/bin/env bash
set -euo pipefail

API_BASE="${API_BASE:-http://localhost:8080}"
USER_ID="${USER_ID:-u1}"
LANG_ID="${LANG_ID:-1}"

need() { command -v "$1" >/dev/null || { echo "Missing $1" >&2; exit 1; }; }
need curl
need jq

echo "Creating packs..."
KITCHEN_ID="$(curl -sS -X POST "$API_BASE/api/packs" \
  -H 'Content-Type: application/json' \
  -d "{\"name\":\"Kitchen\",\"lang_id\":\"$LANG_ID\",\"user_id\":\"$USER_ID\"}" \
  | jq -r '.data.id')"
echo "Kitchen -> $KITCHEN_ID"

ANIMALS_ID="$(curl -sS -X POST "$API_BASE/api/packs" \
  -H 'Content-Type: application/json' \
  -d "{\"name\":\"Animals\",\"lang_id\":\"$LANG_ID\",\"user_id\":\"$USER_ID\"}" \
  | jq -r '.data.id')"
echo "Animals -> $ANIMALS_ID"

tmpdir="$(mktemp -d)"
cleanup() { rm -rf "$tmpdir"; }
trap cleanup EXIT

hindi_name() {
  case "$1" in
    knife) echo "चाकू" ;;
    fork) echo "कांटा" ;;
    spoon) echo "चम्मच" ;;
    plate) echo "थाली" ;;
    cup) echo "कप" ;;
    pan) echo "तवा" ;;
    pot) echo "हांडी" ;;
    stove) echo "चूल्हा" ;;
    fridge) echo "फ्रिज" ;;
    kettle) echo "केतली" ;;
    cat) echo "बिल्ली" ;;
    dog) echo "कुत्ता" ;;
    lion) echo "शेर" ;;
    elephant) echo "हाथी" ;;
    tiger) echo "बाघ" ;;
    bear) echo "भालू" ;;
    horse) echo "घोड़ा" ;;
    bird) echo "पक्षी" ;;
    fish) echo "मछली" ;;
    monkey) echo "बंदर" ;;
    *) echo "$1" ;;
  esac
}

upload_vocab() {
  local pack_id="$1" name="$2" query="$3"
  local file="$tmpdir/${name}"
  local hindi
  hindi="$(hindi_name "$name")"
  # Fetch image from Wikimedia Commons (thumbnail ~800px wide)
  local enc_query url
  enc_query="$(jq -rn --arg q "$query" '$q|@uri')"
  url="$(curl -sS -H 'User-Agent: learnlang-init/1.0' "https://commons.wikimedia.org/w/api.php?action=query&generator=search&gsrsearch=${enc_query}&gsrlimit=1&gsrnamespace=6&prop=imageinfo&iiprop=url&iiurlwidth=800&format=json" \
    | jq -r '(.query.pages[]?.imageinfo[0].thumburl // .query.pages[]?.imageinfo[0].url) // empty' | head -n1)"
  if [[ -z "$url" ]]; then
    # Fallback to a free placeholder image service
    url="https://picsum.photos/800/600?random=1"
  fi
  curl -sSL "$url" -o "$file"
  # Upload via multipart; omit explicit type/extension to let server sniff content and set extension
  curl -sS -X POST "$API_BASE/api/vocabs" \
    -F "name=${name}" \
    -F "translation=${hindi}" \
    -F "pack_id=${pack_id}" \
    -F "image=@${file}" >/dev/null
  echo "Added ${name} (translation: ${hindi})"
}

echo "Seeding Kitchen items..."
upload_vocab "$KITCHEN_ID" knife   "kitchen,knife"
upload_vocab "$KITCHEN_ID" fork    "kitchen,fork,cutlery"
upload_vocab "$KITCHEN_ID" spoon   "kitchen,spoon,cutlery"
upload_vocab "$KITCHEN_ID" plate   "kitchen,plate,dish"
upload_vocab "$KITCHEN_ID" cup     "kitchen,cup,mug"
upload_vocab "$KITCHEN_ID" pan     "kitchen,pan,cookware"
upload_vocab "$KITCHEN_ID" pot     "kitchen,pot,cookware"
upload_vocab "$KITCHEN_ID" stove   "kitchen,stove,oven"
upload_vocab "$KITCHEN_ID" fridge  "kitchen,fridge,refrigerator"
upload_vocab "$KITCHEN_ID" kettle  "kitchen,kettle"

echo "Seeding Animals..."
upload_vocab "$ANIMALS_ID" cat      "animal,cat"
upload_vocab "$ANIMALS_ID" dog      "animal,dog"
upload_vocab "$ANIMALS_ID" lion     "animal,lion"
upload_vocab "$ANIMALS_ID" elephant "animal,elephant"
upload_vocab "$ANIMALS_ID" tiger    "animal,tiger"
upload_vocab "$ANIMALS_ID" bear     "animal,bear"
upload_vocab "$ANIMALS_ID" horse    "animal,horse"
upload_vocab "$ANIMALS_ID" bird     "animal,bird"
upload_vocab "$ANIMALS_ID" fish     "animal,fish"
upload_vocab "$ANIMALS_ID" monkey   "animal,monkey"

echo "Done."