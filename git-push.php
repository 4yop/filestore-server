<?php

while (true)
{
    $time = date("Y-m-d H:i:s",time());
    echo "\n---------start {$time}-----------\n";
    `git pull`;
    `git add -A && git commit -m 'study-go' && git push`;
    echo "\n---------end {$time}-----------\n";
    sleep(600);
}
