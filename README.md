# TSPPD Decision Diagram Code

This source code accompanies the paper:

* Decision Diagrams for Solving Traveling Salesman Problems with Pickup and Delivery in Real Time

All source is licensed under the ZIB Academic License, and can only be used
or referenced for academic purposes. There is no warranty and your mileage
may vary. See the file `LICENSE`.

Building `tsppd-dd` requires `go`, and has only been tested on Linux. To build
it, run `go build`. If that succeeds, you should be able to specify an input file and a diagram form.

```
./tsppd-dd -input <input json file> -form <form>
```

If `-width` is not specified, the resulting diagram will be exact. Otherwise that width controls the restriction and relaxation diagram width. Relaxation and inference duals are specified using the `-relax` and `-infer` flags, respectively. For instance:

```
./tsppd-dd \
    -input grubhub-09-4.json \
    -form sequential \
    -width 5 \
    -infer ap \
    -batch 10 \
    -verbosity 1
```

```
instance        size   form        infer  relax  ordering  width     batch  workers  clock    cpu      primal    optimal  nodes     fails
======================================================================================================================================================
grubhub-09-4    20     sequential  ap     none             5         10     1        0.001    0.001    8169      false    1         0
grubhub-09-4    20     sequential  ap     none             5         10     1        0.004    0.004    7543      false    2         0
grubhub-09-4    20     sequential  ap     none             5         10     1        0.013    0.014    7498      false    42        1
grubhub-09-4    20     sequential  ap     none             5         10     1        0.014    0.015    7429      false    51        2
grubhub-09-4    20     sequential  ap     none             5         10     1        0.016    0.017    7156      false    76        8
grubhub-09-4    20     sequential  ap     none             5         10     1        0.039    0.043    7078      false    270       581
grubhub-09-4    20     sequential  ap     none             5         10     1        0.117    0.132    7078      true     1055      2379
```
