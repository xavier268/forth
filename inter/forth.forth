
( this file redefines some additionnal forth words )

( =============================================== )
( Compiling forth.forth definitions ) 

: version 1 2 ;     ( -- major minor )     
                    ( marker for preloaded  FORTH definitions )

: .version          ( -- display version major. minor )
    ." Forth version " version swap . ." ." . cr
    ;


: +!            ( n addr -- ) ( add n to the cell pointed to by addr )
    dup @       ( n adr x )   
    rot         ( adr x n )   
    + swap !    ( )
    ;        

