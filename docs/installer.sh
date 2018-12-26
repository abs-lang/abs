OS="linux"

if [ "$OSTYPE" = "darwin" ]; then
        OS="darwin"
elif [ "$OSTYPE" = "cygwin"* ]; then
        OS="windows"
elif [ "$OSTYPE" = "win32" ]; then
        OS="windows"
elif [ "$OSTYPE" = "msys" ]; then
        OS="windows"
elif [ "$OSTYPE" = "freebsd" ]; then
        OS="freebds"
fi

ARCH="386"
MACHINE_TYPE=`uname -m`
if [ ${MACHINE_TYPE} = 'x86_64' ]; then
  ARCH="amd64"
fi

echo "Trying to detect the details of your architecture."
echo ""
echo "If these don't seem correct, head over to https://github.com/abs-lang/abs/releases"
echo "and download the right binary for your architecture."
echo ""
echo "OS: $OS"
echo "ARCH: $ARCH"
echo ""
echo "Are these correct? [y/N]"

while read line
do
  INPUT=$(echo $line | awk '{print toupper($0)}')
  if [ $INPUT = "Y" ]; then
    break
  fi
  echo Exited
  exit 1
done < "/dev/stdin"

INSTALLER_PATH=$(dirname $(mktemp -u))
BIN=abs-preview-2-$OS-amd64
cd $INSTALLER_PATH && \
wget https://github.com/abs-lang/abs/releases/download/preview-2/$BIN && \
chmod +x $BIN && \
mv $BIN /usr/local/bin/abs && \
echo "installation completed"
