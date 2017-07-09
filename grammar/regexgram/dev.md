## Dev Notes

### Generate Reducer

A.remove
A.promote
A.promoteValue
A.promoteChildren(1,2,3)

and so forth.

So the entire json spec becomes something like:

Value         -> (string | number | bool | null | Array | Object).promote
Array         -> lb.remove ( Value MoreVals* )? rb.remove
MoreVals      -> comma.remove Value
Object        -> lcb.remove ( KeyVal MoreKeyVals* )? rcb.remove
MoreKeyVals   -> comma.remove KeyVal.promote
KeyVal        -> string.PromoteValue colon.remove Value

### Partial Compiler

It would be nice to be able to output the computed grammar and the reducer. 