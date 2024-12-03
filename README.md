# central-auth

docker build -t central-auth .

docker run -d --network=host -v "$PWD/params/.env:/app/params/.env:Z" -v "$PWD/assets/upload/foto_profil:/app/assets/upload/podcast:Z" -v "$PWD/assets/upload/foto_profil:/app/assets/upload/show:Z" central-auth

Pr1m3@2023
