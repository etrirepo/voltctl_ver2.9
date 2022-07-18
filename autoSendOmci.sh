

while read line; do
  echo $line
  echo `./bossctl device send omcidata $1 $2 $line` > ${line}_history
done < $3

