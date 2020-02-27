( this file redefine some additionnal forth words )

( =============================================== )
." Compiling forth.forth definitions version 0.2 " 

: forth.forth 0 2 ;       ( -- major minor )     
                          ( marker for pre loaded definitions )

      

: VARIABLE  ( -- )        ( VARIABLE <XXX> : creates a variable for XXX ) 
            ( -- addr)    ( Upon exection of XXX )
                          ( Initial variable value is 0 )
    
    HERE 3 +    ( where to store the value )    
    CONSTANT    
    1 ALLOT     ( reserve memory for value )
    ;  

: DECIMAL       10 BASE ! ;
: HEX           16 BASE ! ;  

: +!            ( n addr -- ) ( add n to the cell pointed to by addr )
    DUP @       ( n adr x )   
    ROT         ( adr x n )   
    + SWAP !    ( )
    ;        

