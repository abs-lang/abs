# Simple program that fetches your IP and sums it up
RES=`curl -s 'https://api.ipify.org?format=json' || "ERR"`

if [ "$RES" = "ERR" ]; then
    echo "An error occurred"
    exit 1
fi

IP=`echo $RES | jq -r ".ip"`
IFS=. read first second third fourth <<EOF
${IP##*-}
EOF

total=$((first + second + third + fourth))
if [ $total -gt 100 ]; then
    echo "The sum of [$IP] is a large number, $total."
fi
