( this file redefine some additionnal forth words )

( =============================================== )

." Compiling forth.forth definitions version 1.0 "

: forth.forth 1 0 ;       ( -- major minor )     
                          ( marker for pre loaded definitions )

      

: VARIABLE  ( -- )        ( VARIABLE <XXX> : creates a variable for XXX ) 
            ( -- addr)    ( Upon exection of XXX )
                          ( Initial variable value is 0 )
    
    HERE 3 +    ( where to store the value )    
    CONSTANT    
    1 ALLOT     ( reserve memory for value )
    ;    
    

