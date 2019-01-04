package util

func SliceEqual(a []interface{}, b []interface{}, cp func(interface{}, interface{})bool) bool {
    if len(a) != len(b) {
        return false
    }
    l := len(a)
    for i := 0; i < l; i++ {
        if !cp(a[i], b[i]) {
            return false
        }
    }
    return true
}