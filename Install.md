go build -o mods .
cp ./mods /opt/homebrew/bin/mods-fixed-pipe

cd /Users/mobe/Library/Application Support/mods
ln -sf /Users/mobe/.config/mods/mods.yml
