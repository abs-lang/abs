# Simple program that fetches your IP and sums it up
IP=$(curl -s 'https://api.ipify.org?format=json' | jq -r ".ip")
echo "Your IP is: $IP"
IFS=. read first second third fourth <<EOF
${IP##*-}
EOF
echo "Splitting IP into 4 numbers:"
echo $first
echo $second
echo $third
echo $fourth
echo "Sum:"
total=$((first + second + third + fourth))
echo $total
if [ $total -gt 100 ]
    then
    echo Hey that\'s a large number.
fi