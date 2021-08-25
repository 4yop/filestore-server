<?php

while (true)
{
    echo "\n---------start-----------\n";
    `git pull`;
    `git add -A && git commit -m 'study-go' && git push`;
    echo "\n---------end-----------\n";
    sleep(600);
}
