# central-auth

docker build --network host -t central-auth .

docker run -d --network=host -v "$PWD/params/.env:/app/params/.env:Z" central-auth
