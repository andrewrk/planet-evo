# Planet Evo

Evolution simulation software

## DNA & Reproduction

In this project, DNA is an array of bytes which have
context-based meaning. When a character is interpreted, if the value is
out of range, it is wrapped around (using modulus for example) as many times as
necessary until a value in range is found.

Each cell (particle) of an organism has a complete intact listing of its
DNA, as well as a "currently executing" DNA. There is a pointer to the current
index in the currently executing DNA which henceforth will be named "PC".
The "currently executing" DNA may be modified in place during the course of 
DNA execution.

When cell division occurs, the original cell is left intact, and a new cell
forms in one of the cardinal or diagonal directions. The new cell inherits
the "currently executing" DNA along with the PC as well as the intact listing.
Now both cells will continue executing their now independent DNA programs
at the same position. Some state is different for each cell, however.
For example, cell age is 0 for the new cell and nonzero for the old cell.

When cell death occurs, all stored energy is released, and the cell stops
executing DNA code.

There are 2 registers, X and Y. Instead of parameters being data following the
instructions, all parameters are global variables and each have their own
instruction to set.

* 0 - based on the particle to the (direction), jump to a (label)
* 1 - based on the cell age, jump to a (label)
* 2 - based on the organism age, jump to a (label)
* 3 - wait for cell age modulo (number)
* 4 - wait for organism age modulo (number)
* 5 - wait for organism age (number)
* 6 - set mode to ignore when instruction cannot be run
* 7 - set mode to block when instruction cannot be run - this is the default
* 8 - wait for cell age (number)
* 9 - perform cell death
* A - perform cell division (amount of energy to give new cell, cardinal/diagonal direction, new cell type)
* B - noop
* C - add or subtract to PC (number, positive/negative direction)
* D - ignore the previous instruction
* E - ignore the next instruction
* F - update (register) with (byte)
* G - update (register) with (value of particle to) (direction)
* H - update (register) with (boolean whether or not particle to (direction) is of same organism)
* I - update (register) with (boolean whether or not particle to (direction) is of same cell type)
* J - update (register) with (value of label) * 2 
* K - update (register) with (value of label) + 1
* L - update (register) with (value of label) + 2
* M - update (register) with (my own cell type)
* N - update (label) with (register)
* O - if register modulo 0, goto (label)
* P - if register X is > register Y goto (label)
* Q - update (register) with cell age
* R - update (register) with organism age
* S - perform cell division with fork (new cell energy) (direction) (new cell type) (label for new cell)
* T - perform cell division with variable cell type (new cell energy) (direction) (register)
* U - perform cell division with variable direction (new cell energy) (register) (new cell type)
* V - perform cell division with fork and variable direction (new cell energy) (register) (new cell type) (label for new cell)
* W - update (register1) with (register2) + 1

Let's try that again.

* 0 - Perform cell division. Parameters: CellDivisionEnergyForNewCell, CellDivisionDirection, CellDivisionNewCellType, CellDivisionForkLabel
* 1 - Perform cell death.
* 2 - Jump. Parameters: JumpOperandLeft JumpOperandRight JumpOperation JumpLabel
* 3 - Noop.
* 4 - Wait. Parameters: WaitSource

### Parameter Values

All parameter values have a size of 256. If a value exceeds the bounds, it wraps.

#### JumpOperandLeft / JumpOperandRight / WaitSource

See Value Source

#### Value Source

* 0 - None. This operation is a noop.
* 1 - Register X.
* 2 - Register Y.
* 3 - Cell age.
* 4 - Organism age.
* 5 - The value 0.
* 6 - The value 1.
* 7 - The value 2.
* 8 - Type of particle to the left.
* 9 - Type of particle to the right.
* 10 - Type of particle to the up.
* 11 - Type of particle to the down.
* 12 - Type of particle to the top/left.
* 13 - Type of particle to the top/right.
* 14 - Type of particle to the bottom/left.
* 15 - Type of particle to the bottom/right.
* 16 - Count of cells in organism.

#### JumpLabel

See Code Label

#### CellDivisionEnergyForNewCell

Number, 0 - 100

#### CellDivisionDirection

See Cardinal/Diagonal Direction

#### CellDivisionNewCellType

See New Cell Type

#### CellDivisionForkLabel

See Code Label

#### Code Label

Number, 0 - 256. Number of bytes offset. When a mutation occurs, this number
is adjusted to point to the same location.

#### Cardinal/Diagonal Direction
* 0 - left
* 1 - right
* 2 - up
* 3 - down
* 4 - top/left
* 5 - top/right
* 6 - bottom/left
* 7 - bottom/right

#### Positive/Negative Direction
* 0 - positive
* 1 - negative

#### New Cell Type
* 0 - Chloro
* 1 - Fiber
