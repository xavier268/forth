( this file redefine some additionnal forth words )

( =============================================== )
( Compiling forth.forth definitions ) 

: VERSION 0 3 ;     ( -- major minor )     
                    ( marker for preloaded  FORTH definitions )

: CONSTANT  ( value -- )
<BUILDS ( value -- ) ( constant value should be available on DATA stack )
, ( ) 
DOES> ( addr -- ... )
@ 
( ... -- value )
;

: VARIABLE 
 <BUILDS ( -- )
  1 ALLOT ( allot 1 data cell )
  DOES> ( dataAddr -- ... )
  ( fine, data addr already on stack, do nothing ! )
  ( ... -- dataAddr )
;

: DECIMAL       10 BASE ! ;
: HEX           16 BASE ! ;  

: +!            ( n addr -- ) ( add n to the cell pointed to by addr )
    DUP @       ( n adr x )   
    ROT         ( adr x n )   
    + SWAP !    ( )
    ;        

