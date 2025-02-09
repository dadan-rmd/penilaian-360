# penilaian-360

docker build --network host -t penilaian-360 .

docker run -d --network=host -v "$PWD/params/.env:/app/params/.env:Z" penilaian-360
