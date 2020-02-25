( this file redefine some additionnal forth words )

: VARIABLE  ( -- )        ( VARIABLE <XXX> : creates a variable for XXX ) 
            ( -- addr)    ( Upon exection of XXX )
    
    HERE 3 +    ( where to store the value )    
    CONSTANT    
    1 ALLOT     ( reserve memory for value )
    ;    
    