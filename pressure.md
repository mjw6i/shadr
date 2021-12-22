# pressure
there is one dimensional array of integers, representing pressure of a given tile\
on each iteration pressure will normalise following those rules:\
 - pressure transfer happens when checking tile with higher pressure
 - pressure is transfered with tile on the left before tile on the right
 - maximum pressure that can be transfered in one iteration is `floor(higer_pressure / 20)`
 - if the pressure is not equal after the transfer, tile that had higher pressure should still have higher pressure
