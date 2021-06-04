<?php

# Client code

echo "Starting\n";

# Create our client object.
$gmclient= new GearmanClient();

# Add default server ();
$gmclient->addServer();

$gmclient->setCreatedCallback("reverse_created");
$gmclient->setDataCallback("reverse_data");
$gmclient->setStatusCallback("reverse_status");
$gmclient->setCompleteCallback("reverse_complete");
$gmclient->setFailCallback("reverse_fail");

# set some arbitrary application data

# add two tasks
  $task1= $gmclient->doBackGround("reverse","What");
  $task2= $gmclient->doBackGround("reverse","Where");


function reverse_created($task)
{
   echo "CREATED: ". $task->jobHandle()."\n";
}

function reverse_status($task){

   echo "STATUS: ".$task->jobHandle()."-".$task->taskNumerator()."/".$task->taskDenominator()."\n";
    
}

function reverse_complete($task)
{
    echo "COMPLETE: " . $task->jobHandle() . ", " . $task->data() . "\n";
    

    
}

function reverse_fail($task)
{
    echo "FAILED: " . $task->jobHandle() . "\n";
}

function reverse_data($task)
{
    echo "DATA: " . $task->data() . "\n";
     
}






?>



