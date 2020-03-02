( this file redefine some additionnal forth words )

( =============================================== )
." Compiling forth.forth definitions version 0.3 " 

: forth.forth 0 3 ;       ( -- major minor )     
                          ( marker for pre loaded definitions )

: CONSTANT  ( value -- )
<BUILDS ( value -- )
1 ALLOT HERE 1 - ! ( ) 
DOES> ( addr -- )
@ ( -- value )
;

: VARIABLE 
 <BUILDS 
  1 ALLOT ( allot 1 data cell )
  DOES> ( dataAddr -- ) ( fine, do nothing ! )
;

: DECIMAL       10 BASE ! ;
: HEX           16 BASE ! ;  

: +!            ( n addr -- ) ( add n to the cell pointed to by addr )
    DUP @       ( n adr x )   
    ROT         ( adr x n )   
    + SWAP !    ( )
    ;        

