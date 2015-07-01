// auto generated by go tool dist

package runtime
import "unsafe"
var _ unsafe.Pointer

const _Gidle = 0
const _Grunnable = 1
const _Grunning = 2
const _Gsyscall = 3
const _Gwaiting = 4
const _Gmoribund_unused = 5
const _Gdead = 6
const _Genqueue = 7
const _Gcopystack = 8
const _Gscan = 4096
const _Gscanrunnable = 4097
const _Gscanrunning = 4098
const _Gscansyscall = 4099
const _Gscanwaiting = 4100
const _Gscanenqueue = 4103
const _Pidle = 0
const _Prunning = 1
const _Psyscall = 2
const _Pgcstop = 3
const _Pdead = 4
const _PtrSize = 4
type mutex struct {
	key	uintptr
}

type note struct {
	key	uintptr
}

type _string struct {
	str	*byte
	len	int
}

type funcval struct {
	fn	unsafe.Pointer
}

type iface struct {
	tab	*itab
	data	unsafe.Pointer
}

type eface struct {
	_type	*_type
	data	unsafe.Pointer
}

type _complex64 struct {
	real	float32
	imag	float32
}

type _complex128 struct {
	real	float64
	imag	float64
}

type slice struct {
	array	*byte
	len	uint
	cap	uint
}

type gobuf struct {
	sp	uintptr
	pc	uintptr
	g	*g
	ctxt	unsafe.Pointer
	ret	uintreg
	lr	uintptr
}

type sudog struct {
	g	*g
	selectdone	*uint32
	next	*sudog
	prev	*sudog
	elem	unsafe.Pointer
	releasetime	int64
	nrelease	int32
	waitlink	*sudog
}

type gcstats struct {
	nhandoff	uint64
	nhandoffcnt	uint64
	nprocyield	uint64
	nosyield	uint64
	nsleep	uint64
}

type libcall struct {
	fn	uintptr
	n	uintptr
	args	uintptr
	r1	uintptr
	r2	uintptr
	err	uintptr
}

type wincallbackcontext struct {
	gobody	unsafe.Pointer
	argsize	uintptr
	restorestack	uintptr
	cleanstack	bool
}

type stack struct {
	lo	uintptr
	hi	uintptr
}

type g struct {
	stack	stack
	stackguard0	uintptr
	stackguard1	uintptr
	_panic	*_panic
	_defer	*_defer
	sched	gobuf
	syscallsp	uintptr
	syscallpc	uintptr
	param	unsafe.Pointer
	atomicstatus	uint32
	goid	int64
	waitsince	int64
	waitreason	string
	schedlink	*g
	issystem	bool
	preempt	bool
	paniconfault	bool
	preemptscan	bool
	gcworkdone	bool
	throwsplit	bool
	raceignore	int8
	m	*m
	lockedm	*m
	sig	int32
	writebuf	[]byte
	sigcode0	uintptr
	sigcode1	uintptr
	sigpc	uintptr
	gopc	uintptr
	racectx	uintptr
	waiting	*sudog
	end	[0]uintptr
}

type m struct {
	g0	*g
	morebuf	gobuf
	procid	uint64
	gsignal	*g
	tls	[4]uintptr
	mstartfn	unsafe.Pointer
	curg	*g
	caughtsig	*g
	p	*p
	nextp	*p
	id	int32
	mallocing	int32
	throwing	int32
	gcing	int32
	locks	int32
	softfloat	int32
	dying	int32
	profilehz	int32
	helpgc	int32
	spinning	bool
	blocked	bool
	fastrand	uint32
	ncgocall	uint64
	ncgo	int32
	cgomal	*cgomal
	park	note
	alllink	*m
	schedlink	*m
	machport	uint32
	mcache	*mcache
	lockedg	*g
	createstack	[32]uintptr
	freglo	[16]uint32
	freghi	[16]uint32
	fflag	uint32
	locked	uint32
	nextwaitm	*m
	waitsema	uintptr
	waitsemacount	uint32
	waitsemalock	uint32
	gcstats	gcstats
	needextram	bool
	traceback	uint8
	waitunlockf	unsafe.Pointer
	waitlock	unsafe.Pointer
	scalararg	[4]uintptr
	ptrarg	[4]unsafe.Pointer
	end	[0]uintptr
}

type p struct {
	lock	mutex
	id	int32
	status	uint32
	link	*p
	schedtick	uint32
	syscalltick	uint32
	m	*m
	mcache	*mcache
	deferpool	[5]*_defer
	goidcache	uint64
	goidcacheend	uint64
	runqhead	uint32
	runqtail	uint32
	runq	[256]*g
	gfree	*g
	gfreecnt	int32
	pad	[64]byte
}

const _MaxGomaxprocs = 256
type schedt struct {
	lock	mutex
	goidgen	uint64
	midle	*m
	nmidle	int32
	nmidlelocked	int32
	mcount	int32
	maxmcount	int32
	pidle	*p
	npidle	uint32
	nmspinning	uint32
	runqhead	*g
	runqtail	*g
	runqsize	int32
	gflock	mutex
	gfree	*g
	ngfree	int32
	gcwaiting	uint32
	stopwait	int32
	stopnote	note
	sysmonwait	uint32
	sysmonnote	note
	lastpoll	uint64
	profilehz	int32
}

const _LockExternal = 1
const _LockInternal = 2
type sigtab struct {
	flags	int32
	name	*int8
}

const _SigNotify = 1
const _SigKill = 2
const _SigThrow = 4
const _SigPanic = 8
const _SigDefault = 16
const _SigHandling = 32
const _SigIgnored = 64
const _SigGoExit = 128
type _func struct {
	entry	uintptr
	nameoff	int32
	args	int32
	frame	int32
	pcsp	int32
	pcfile	int32
	pcln	int32
	npcdata	int32
	nfuncdata	int32
}

type itab struct {
	inter	*interfacetype
	_type	*_type
	link	*itab
	bad	int32
	unused	int32
	fun	[0]unsafe.Pointer
}

const _NaCl = 0
const _Windows = 0
const _Solaris = 0
const _Plan9 = 0
type lfnode struct {
	next	*lfnode
	pushcnt	uintptr
}

type parfor struct {
	body	unsafe.Pointer
	done	uint32
	nthr	uint32
	nthrmax	uint32
	thrseq	uint32
	cnt	uint32
	ctx	unsafe.Pointer
	wait	bool
	thr	*parforthread
	pad	uint32
	nsteal	uint64
	nstealcnt	uint64
	nprocyield	uint64
	nosyield	uint64
	nsleep	uint64
}

type cgomal struct {
	next	*cgomal
	alloc	unsafe.Pointer
}

type debugvars struct {
	allocfreetrace	int32
	efence	int32
	gctrace	int32
	gcdead	int32
	scheddetail	int32
	schedtrace	int32
	scavenge	int32
}

const _GCoff = 0
const _GCquiesce = 1
const _GCstw = 2
const _GCmark = 3
const _GCsweep = 4
type forcegcstate struct {
	lock	mutex
	g	*g
	idle	uint32
}

var gcphase	uint32
const _Structrnd = 4
var startup_random_data	*byte
var startup_random_data_len	uint32
var invalidptr	int32
const _HashRandomBytes = 32
type _defer struct {
	siz	int32
	started	bool
	argp	uintptr
	pc	uintptr
	fn	*funcval
	_panic	*_panic
	link	*_defer
}

type _panic struct {
	argp	unsafe.Pointer
	arg	interface{}
	link	*_panic
	recovered	bool
	aborted	bool
}

type stkframe struct {
	fn	*_func
	pc	uintptr
	continpc	uintptr
	lr	uintptr
	sp	uintptr
	fp	uintptr
	varp	uintptr
	argp	uintptr
	arglen	uintptr
	argmap	*bitvector
}

const _TraceRuntimeFrames = 1
const _TraceTrap = 2
const _TracebackMaxFrames = 100
var emptystring	string
var allg	**g
var allglen	uintptr
var lastg	*g
var allm	*m
var allp	[257]*p
var gomaxprocs	int32
var needextram	uint32
var panicking	uint32
var goos	*int8
var ncpu	int32
var iscgo	bool
var sysargs	unsafe.Pointer
var maxstring	uintptr
var cpuid_ecx	uint32
var cpuid_edx	uint32
var debug	debugvars
var maxstacksize	uintptr
var signote	note
var forcegc	forcegcstate
var sched	schedt
var newprocs	int32
var worldsema	uint32
var nan	float64
var posinf	float64
var neginf	float64
const _UseSpanType = 1
const thechar = 53
const _BigEndian = 0
const _CacheLineSize = 32
const _RuntimeGogoBytes = 60
const _PhysPageSize = 4096
const _PCQuantum = 4
const _Int64Align = 4
const _PageShift = 13
const _PageSize = 8192
const _PageMask = 8191
const _NumSizeClasses = 67
const _MaxSmallSize = 32768
const _TinySize = 16
const _TinySizeClass = 2
const _FixAllocChunk = 16384
const _MaxMHeapList = 128
const _HeapAllocChunk = 1048576
const _StackCacheSize = 32768
const _NumStackOrders = 3
const _MHeapMap_Bits = 19
const _MaxGcproc = 32
type mlink struct {
	next	*mlink
}

type fixalloc struct {
	size	uintptr
	first	unsafe.Pointer
	arg	unsafe.Pointer
	list	*mlink
	chunk	*byte
	nchunk	uint32
	inuse	uintptr
	stat	*uint64
}

type mstatsbysize struct {
	size	uint32
	nmalloc	uint64
	nfree	uint64
}

type mstats struct {
	alloc	uint64
	total_alloc	uint64
	sys	uint64
	nlookup	uint64
	nmalloc	uint64
	nfree	uint64
	heap_alloc	uint64
	heap_sys	uint64
	heap_idle	uint64
	heap_inuse	uint64
	heap_released	uint64
	heap_objects	uint64
	stacks_inuse	uint64
	stacks_sys	uint64
	mspan_inuse	uint64
	mspan_sys	uint64
	mcache_inuse	uint64
	mcache_sys	uint64
	buckhash_sys	uint64
	gc_sys	uint64
	other_sys	uint64
	next_gc	uint64
	last_gc	uint64
	pause_total_ns	uint64
	pause_ns	[256]uint64
	pause_end	[256]uint64
	numgc	uint32
	enablegc	bool
	debuggc	bool
	by_size	[67]mstatsbysize
	tinyallocs	uint64
}

var memstats	mstats
var class_to_size	[67]int32
var class_to_allocnpages	[67]int32
var size_to_class8	[129]int8
var size_to_class128	[249]int8
type mcachelist struct {
	list	*mlink
	nlist	uint32
}

type stackfreelist struct {
	list	*mlink
	size	uintptr
}

type mcache struct {
	next_sample	int32
	local_cachealloc	intptr
	tiny	*byte
	tinysize	uintptr
	local_tinyallocs	uintptr
	alloc	[67]*mspan
	stackcache	[3]stackfreelist
	sudogcache	*sudog
	gcworkbuf	unsafe.Pointer
	local_nlookup	uintptr
	local_largefree	uintptr
	local_nlargefree	uintptr
	local_nsmallfree	[67]uintptr
}

const _KindSpecialFinalizer = 1
const _KindSpecialProfile = 2
type special struct {
	next	*special
	offset	uint16
	kind	byte
}

type specialfinalizer struct {
	special	special
	fn	*funcval
	nret	uintptr
	fint	*_type
	ot	*ptrtype
}

type specialprofile struct {
	special	special
	b	*bucket
}

const _MSpanInUse = 0
const _MSpanStack = 1
const _MSpanFree = 2
const _MSpanListHead = 3
const _MSpanDead = 4
type mspan struct {
	next	*mspan
	prev	*mspan
	start	pageID
	npages	uintptr
	freelist	*mlink
	sweepgen	uint32
	ref	uint16
	sizeclass	uint8
	incache	bool
	state	uint8
	needzero	uint8
	elemsize	uintptr
	unusedsince	int64
	npreleased	uintptr
	limit	*byte
	speciallock	mutex
	specials	*special
}

type mcentral struct {
	lock	mutex
	sizeclass	int32
	nonempty	mspan
	empty	mspan
}

type mheapcentral struct {
	mcentral	mcentral
	pad	[32]byte
}

type mheap struct {
	lock	mutex
	free	[128]mspan
	freelarge	mspan
	busy	[128]mspan
	busylarge	mspan
	allspans	**mspan
	gcspans	**mspan
	nspan	uint32
	nspancap	uint32
	sweepgen	uint32
	sweepdone	uint32
	spans	**mspan
	spans_mapped	uintptr
	bitmap	*byte
	bitmap_mapped	uintptr
	arena_start	*byte
	arena_used	*byte
	arena_end	*byte
	arena_reserved	bool
	central	[67]mheapcentral
	spanalloc	fixalloc
	cachealloc	fixalloc
	specialfinalizeralloc	fixalloc
	specialprofilealloc	fixalloc
	speciallock	mutex
	largefree	uint64
	nlargefree	uint64
	nsmallfree	[67]uint64
}

var mheap_	mheap
var gcpercent	int32
const _FlagNoScan = 1
const _FlagNoZero = 2
type finalizer struct {
	fn	*funcval
	arg	unsafe.Pointer
	nret	uintptr
	fint	*_type
	ot	*ptrtype
}

type finblock struct {
	alllink	*finblock
	next	*finblock
	cnt	int32
	cap	int32
	fin	[1]finalizer
}

var finlock	mutex
var fing	*g
var fingwait	bool
var fingwake	bool
var finq	*finblock
var finc	*finblock
type bitvector struct {
	n	int32
	bytedata	*uint8
}

type stackmap struct {
	n	int32
	nbit	int32
	bytedata	[0]uint8
}

var gcdatamask	bitvector
var gcbssmask	bitvector
type _type struct {
	size	uintptr
	hash	uint32
	_unused	uint8
	align	uint8
	fieldalign	uint8
	kind	uint8
	alg	unsafe.Pointer
	gc	[2]uintptr
	_string	*string
	x	*uncommontype
	ptrto	*_type
	zero	*byte
}

type method struct {
	name	*string
	pkgpath	*string
	mtyp	*_type
	typ	*_type
	ifn	unsafe.Pointer
	tfn	unsafe.Pointer
}

type uncommontype struct {
	name	*string
	pkgpath	*string
	mhdr	[]byte
	m	[0]method
}

type imethod struct {
	name	*string
	pkgpath	*string
	_type	*_type
}

type interfacetype struct {
	typ	_type
	mhdr	[]byte
	m	[0]imethod
}

type maptype struct {
	typ	_type
	key	*_type
	elem	*_type
	bucket	*_type
	hmap	*_type
	keysize	uint8
	indirectkey	bool
	valuesize	uint8
	indirectvalue	bool
	bucketsize	uint16
}

type chantype struct {
	typ	_type
	elem	*_type
	dir	uintptr
}

type slicetype struct {
	typ	_type
	elem	*_type
}

type functype struct {
	typ	_type
	dotdotdot	bool
	in	[]byte
	out	[]byte
}

type ptrtype struct {
	typ	_type
	elem	*_type
}

type waitq struct {
	first	*sudog
	last	*sudog
}

type hchan struct {
	qcount	uint
	dataqsiz	uint
	buf	*byte
	elemsize	uint16
	closed	uint32
	elemtype	*_type
	sendx	uint
	recvx	uint
	recvq	waitq
	sendq	waitq
	lock	mutex
}

const _CaseRecv = 1
const _CaseSend = 2
const _CaseDefault = 3
type scase struct {
	elem	unsafe.Pointer
	_chan	*hchan
	pc	uintptr
	kind	uint16
	so	uint16
	receivedp	*bool
	releasetime	int64
}

type _select struct {
	tcase	uint16
	ncase	uint16
	pollorder	*uint16
	lockorder	**hchan
	scase	[1]scase
}

const _EINTR = 4
const _EFAULT = 14
const _PROT_NONE = 0
const _PROT_READ = 1
const _PROT_WRITE = 2
const _PROT_EXEC = 4
const _MAP_ANON = 4096
const _MAP_PRIVATE = 2
const _MAP_FIXED = 16
const _MADV_FREE = 6
const _SA_SIGINFO = 64
const _SA_RESTART = 2
const _SA_ONSTACK = 1
const _SIGHUP = 1
const _SIGINT = 2
const _SIGQUIT = 3
const _SIGILL = 4
const _SIGTRAP = 5
const _SIGABRT = 6
const _SIGEMT = 7
const _SIGFPE = 8
const _SIGKILL = 9
const _SIGBUS = 10
const _SIGSEGV = 11
const _SIGSYS = 12
const _SIGPIPE = 13
const _SIGALRM = 14
const _SIGTERM = 15
const _SIGURG = 16
const _SIGSTOP = 17
const _SIGTSTP = 18
const _SIGCONT = 19
const _SIGCHLD = 20
const _SIGTTIN = 21
const _SIGTTOU = 22
const _SIGIO = 23
const _SIGXCPU = 24
const _SIGXFSZ = 25
const _SIGVTALRM = 26
const _SIGPROF = 27
const _SIGWINCH = 28
const _SIGINFO = 29
const _SIGUSR1 = 30
const _SIGUSR2 = 31
const _FPE_INTDIV = 1
const _FPE_INTOVF = 2
const _FPE_FLTDIV = 3
const _FPE_FLTOVF = 4
const _FPE_FLTUND = 5
const _FPE_FLTRES = 6
const _FPE_FLTINV = 7
const _FPE_FLTSUB = 8
const _BUS_ADRALN = 1
const _BUS_ADRERR = 2
const _BUS_OBJERR = 3
const _SEGV_MAPERR = 1
const _SEGV_ACCERR = 2
const _ITIMER_REAL = 0
const _ITIMER_VIRTUAL = 1
const _ITIMER_PROF = 2
const _EV_ADD = 1
const _EV_DELETE = 2
const _EV_CLEAR = 32
const _EV_RECEIPT = 0
const _EV_ERROR = 16384
const _EVFILT_READ = 0
const _EVFILT_WRITE = 1
type sigaltstackt struct {
	ss_sp	*byte
	ss_size	uint32
	ss_flags	int32
}

type sigset struct {
	__bits	[4]uint32
}

type siginfo struct {
	_signo	int32
	_code	int32
	_errno	int32
	_reason	[20]byte
}

type stackt struct {
	ss_sp	*byte
	ss_size	uint32
	ss_flags	int32
}

type timespec struct {
	tv_sec	int64
	tv_nsec	int32
}

type timeval struct {
	tv_sec	int64
	tv_usec	int32
}

type itimerval struct {
	it_interval	timeval
	it_value	timeval
}

type mcontextt struct {
	__gregs	[17]uint32
	__fpu	[140]byte
	_mc_tlsbase	uint32
}

type ucontextt struct {
	uc_flags	uint32
	uc_link	*ucontextt
	uc_sigmask	sigset
	uc_stack	stackt
	uc_mcontext	mcontextt
	__uc_pad	[2]int32
}

type keventt struct {
	ident	uint32
	filter	uint32
	flags	uint32
	fflags	uint32
	data	int64
	udata	*byte
}

const _REG_R0 = 0
const _REG_R1 = 1
const _REG_R2 = 2
const _REG_R3 = 3
const _REG_R4 = 4
const _REG_R5 = 5
const _REG_R6 = 6
const _REG_R7 = 7
const _REG_R8 = 8
const _REG_R9 = 9
const _REG_R10 = 10
const _REG_R11 = 11
const _REG_R12 = 12
const _REG_R13 = 13
const _REG_R14 = 14
const _REG_R15 = 15
const _REG_CPSR = 16
const _SS_DISABLE = 4
const _SIG_BLOCK = 1
const _SIG_UNBLOCK = 2
const _SIG_SETMASK = 3
const _NSIG = 33
const _SI_USER = 0
const _UC_SIGMASK = 1
const _UC_CPU = 4


























































const _KindBool = 1
const _KindInt = 2
const _KindInt8 = 3
const _KindInt16 = 4
const _KindInt32 = 5
const _KindInt64 = 6
const _KindUint = 7
const _KindUint8 = 8
const _KindUint16 = 9
const _KindUint32 = 10
const _KindUint64 = 11
const _KindUintptr = 12
const _KindFloat32 = 13
const _KindFloat64 = 14
const _KindComplex64 = 15
const _KindComplex128 = 16
const _KindArray = 17
const _KindChan = 18
const _KindFunc = 19
const _KindInterface = 20
const _KindMap = 21
const _KindPtr = 22
const _KindSlice = 23
const _KindString = 24
const _KindStruct = 25
const _KindUnsafePointer = 26
const _KindDirectIface = 32
const _KindGCProg = 64
const _KindNoPointers = 128
const _KindMask = 31
const _StackSystem = 0
const _StackMin = 2048
const _FixedStack0 = 2048
const _FixedStack1 = 2047
const _FixedStack2 = 2047
const _FixedStack3 = 2047
const _FixedStack4 = 2047
const _FixedStack5 = 2047
const _FixedStack6 = 2047
const _FixedStack = 2048
const _StackBig = 4096
const _StackGuard = 512
const _StackSmall = 128
const _StackLimit = 384
var sizeof_c_mstats	uintptr
var maxmem	uintptr
var end	[0]byte
















































var memprofilerate	int
var emptymspan	mspan
















































const gcBits = 4
const wordsPerBitmapByte = 2
const insData = 1
const insArray = 2
const insArrayEnd = 3
const insEnd = 4
const _BitsPerPointer = 2
const _BitsMask = 3
const _PointersPerByte = 4
const _BitsDead = 0
const _BitsScalar = 1
const _BitsPointer = 2
const _BitsMultiWord = 3
const _BitsIface = 2
const _BitsEface = 3
const _MaxGCMask = 64
const bitBoundary = 1
const bitMarked = 2
const bitMask = 3
const bitPtrMask = 12














const _Debug = 0
const _DebugPtrs = 0
const _ConcurrentSweep = 1
const _WorkbufSize = 4096
const _FinBlockSize = 4096
const _RootData = 0
const _RootBss = 1
const _RootFinalizers = 2
const _RootSpans = 3
const _RootFlushCaches = 4
const _RootCount = 5
var oneptr	[0]byte
type workbuf struct {
	node	lfnode
	nobj	uintptr
	obj	[1021]*byte
}

var data	[0]byte
var edata	[0]byte
var bss	[0]byte
var ebss	[0]byte
var gcdata	[0]byte
var gcbss	[0]byte
var finptrmask	[256]byte
var allfin	*finblock
var gclock	mutex
var badblock	[1024]uintptr
var nbadblock	int32
var bgsweepv	funcval
type workdata struct {
	full	uint64
	empty	uint64
	pad0	[32]byte
	nproc	uint32
	tstart	int64
	nwait	uint32
	ndone	uint32
	alldone	note
	markfor	*parfor
	spans	**mspan
	nspan	uint32
}

var work	workdata
var finalizer1	[0]byte
type sweepdata struct {
	g	*g
	parked	bool
	spanidx	uint32
	nbgsweep	uint32
	npausesweep	uint32
}

var sweep	sweepdata
type gc_args struct {
	start_time	int64
	eagersweep	bool
}

const bitmapChunk = 8192


























































const _GoidCacheBatch = 16
var m0	m
var g0	g
var extram	*m
var allglock	mutex
var buildversion	string
var _cgo_init	unsafe.Pointer
var _cgo_malloc	unsafe.Pointer
var _cgo_free	unsafe.Pointer
var cgomalloc	unsafe.Pointer
var cgofree	unsafe.Pointer
var _cgo_thread_start	unsafe.Pointer
type cgothreadstart struct {
	g	*g
	tls	*uintptr
	fn	unsafe.Pointer
}

type profstate struct {
	lock	uint32
	hz	int32
}

var prof	profstate
var etext	[0]byte
type pdesc struct {
	schedtick	uint32
	schedwhen	int64
	syscalltick	uint32
	syscallwhen	int64
}

var experiment	[0]int8
















































type parforthread struct {
	pos	uint64
	nsteal	uint64
	nstealcnt	uint64
	nprocyield	uint64
	nosyield	uint64
	nsleep	uint64
	pad	[32]byte
}



























































const _StackDebug = 0
const _StackFromSystem = 0
const _StackFaultOnFree = 0
const _StackPoisonCopy = 0
const _StackCache = 1
var stackpool	[3]mspan
var stackpoolmu	mutex
var stackfreequeue	stack
var mapnames	[0]*uint8
type adjustinfo struct {
	old	stack
	delta	uintptr
}

