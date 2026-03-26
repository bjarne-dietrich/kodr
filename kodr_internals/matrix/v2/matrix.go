package matrix

import (
	"runtime"
	"sync"

	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/kodr_internals/operations"
)

type Matrix [][]byte

// Rows returns the current number of rows in Matrix.
func (m *Matrix) Rows() uint {
	return uint(len(*m))
}

// Cols returns the current number of columns in Matrix.
func (m *Matrix) Cols() uint {
	return uint(len((*m)[0]))
}

// Multiply returns the multiplication of Matrix and with.
func (m *Matrix) Multiply(with Matrix) (Matrix, error) {
	mr := int(m.Rows())
	mc := int(m.Cols())
	wr := int(with.Rows())
	wc := int(with.Cols())

	if mc != wr {
		return nil, kodr.ErrMatrixDimensionMismatch
	}

	out := make(Matrix, mr)
	for i := 0; i < mr; i++ {
		out[i] = make([]byte, wc)
	}

	workers := runtime.GOMAXPROCS(0)
	if workers > mr {
		workers = mr
	}
	rowsPerChunk := 4
	minChunks := workers * 2 // keep workers fed
	if workers <= 1 || mr < rowsPerChunk*minChunks {
		for i := 0; i < mr; i++ {
			mulRowKernel((*m)[i], out[i], with, mc)
		}
		return out, nil
	}

	jobs := make(chan int, workers*2)

	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for start := range jobs {
				end := start + rowsPerChunk
				if end > mr {
					end = mr
				}
				for i := start; i < end; i++ {
					mulRowKernel((*m)[i], out[i], with, mc)
				}
			}
		}()
	}

	for start := 0; start < mr; start += rowsPerChunk {
		jobs <- start
	}
	close(jobs)
	wg.Wait()

	return out, nil
}

func mulRowKernel(ai []byte, ci []byte, with Matrix, mc int) {
	for k := 0; k < mc; k++ {
		c := ai[k] // scalar A[i,k]
		// C[i,:] += B[k,:] * c
		operations.MulAddConst(ci, with[k], c)
	}
}

func (m *Matrix) Transposed() Matrix {
	r := int(m.Rows())
	c := int(m.Cols())

	out := make(Matrix, c)
	for j := 0; j < c; j++ {
		row := make([]byte, r)
		for i := 0; i < r; i++ {
			row[i] = (*m)[i][j]
		}
		out[j] = row
	}
	return out
}
