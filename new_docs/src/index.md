---
home: true
heroImage: /abs-horizontal.png
tagline: Home of the ABS programming language - the joy of shell scripting.
actionText: Quick Start →
actionLink: /introduction/
features:
  - title: Feature 1 Title
    details: Feature 1 Description
  - title: Feature 2 Title
    details: Feature 2 Description
  - title: Feature 3 Title
    details: Feature 3 Description
footer: Made by Alex Nadalin with ❤️
---

::: slot sample-code

```sh
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
```

:::
