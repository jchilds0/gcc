{
    int i; int j; float[100] a; float v; float x;
    while ( true ) {
        do i = i+1; while (a[i] < v);
        do j = j-1; while (a[j] > v);
        if ( i >= j ) break;
        x = a[i]; a[i] = a[j]; a[j] = x;
    }
}
}