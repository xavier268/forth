( test file )
CR ." Compiling test.forth file " CR

( Work In Progress : testing build/does ... )


: VARIABLE [ ." Redefining VARIABLE " CR ] 
 <BUILDS 
  ." Creating head of a new VARIABLE " 
  1 ALLOT ( allot 1 data cell )
  DOES> ( dataAddr -- ) ( fine, do nothing ! )
  DUP ." Accessing variable at address : " . CR ( <--- DEBUG TRACE )
;

: CONSTANT [ ." Redefining CONSTANT " CR ] ( value -- )
<BUILDS 
." Creating head of a new CONSTANT "
1 ALLOT
( value -- )
HERE 1 - ! ( ) 
DOES> ( addr -- )
@ ( -- value )
DUP ." read value : " . ( <--- DEBUG TRACE )
;


