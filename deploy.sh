
set -e
mkdir -p /opt/amidesk-api/data
if [ ! -f /opt/amidesk-api/data/server.yaml ]; then
  echo "data/server.yaml does not exist. Copying and configuring..."
  cp /opt/amidesk-api/backend/server.yaml /opt/amidesk-api/data/server.yaml
  NEW_KEY=$(openssl rand -hex 32)
  sed -i "s|^signKey:.*|signKey: \"$NEW_KEY\"|g" /opt/amidesk-api/data/server.yaml
  sed -i "s|debugMode: true|debugMode: false|g" /opt/amidesk-api/data/server.yaml
  echo "Configuration complete."
else
  echo "data/server.yaml already exists. Skipping configuration."
fi

