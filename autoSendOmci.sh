

while read line; do
  echo $line
  echo `./bossctl device send omcidata $1 $2 $line`
done < $3

