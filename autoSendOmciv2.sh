

while read line; do
  echo $line
  echo `./bossctl device send omcidatav2 $1 $2 $3 $line`
done < $4

