<?php

while (true)
{
    $time = time();
    echo "\n---------start-----------\n";
    `git pull`;
    `git add -A && git commit -m 'study-go' && git push`;
    echo "\n---------end {$time}-----------\n";
    sleep(600);
}
