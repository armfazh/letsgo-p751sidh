package cln16sidh

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

// Sage script for generating test vectors:
// sage: p = 2^372 * 3^239 - 1; Fp = GF(p)
// sage: R.<x> = Fp[]
// sage: Fp2 = Fp.extension(x^2 + 1, 'i')
// sage: i = Fp2.gen()
// sage: E = EllipticCurve(Fp2, [0,A/C,0,1,0])
// sage: X, Y, Z = (817215127176107155479622194880146209497224298781185275314486552489943358359683935722341108891938834236465163218045208196051151604093542873782962420642628777425511 4241789158000915683252363913079335550843837650671094705509470594*i + 93265748580399441216040154393817201955561834227195054974485410732727205450477422355269637733590040218389 61919129020087515274115525812121436661025030481584576474033630899768377131534320053412545346268645085054880212827284581557, 2381174772709336084066332457520782192315178511983342038392622832616744048226360647551642232950959910067260611740876401494529727990031260499974773548012283808741733925525689 114517493995359390158666069816204787133942283380884077*i + 537895623203422833518969796914455655278385875583228419480247092297605464569632411896633315826744276713852822796884 1257817537239745277092206433048875637709652271370008564179304718555812947398374153513738054572355903547642836171, 1)
// sage: P = E((X,Y,Z))
// sage: X2, Y2, Z2 = 2*P
// sage: X3, Y3, Z3 = 3*P

// A = 4385300808024233870220415655826946795549183378139271271040522089756750951667981765872679172832050962894122367066234419550072004266298327417513857609747116903999863022476533671840646615759860564818837299058134292387429068536219*i + 1408083354499944307008104531475821995920666351413327060806684084512082259107262519686546161682384352696826343970108773343853651664489352092568012759783386151707999371397181344707721407830640876552312524779901115054295865393760
var curve_A = ExtensionFieldElement{a: fp751Element{0x8319eb18ca2c435e, 0x3a93beae72cd0267, 0x5e465e1f72fd5a84, 0x8617fa4150aa7272, 0x887da24799d62a13, 0xb079b31b3c7667fe, 0xc4661b150fa14f2e, 0xd4d2b2967bc6efd6, 0x854215a8b7239003, 0x61c5302ccba656c2, 0xf93194a27d6f97a2, 0x1ed9532bca75}, b: fp751Element{0xb6f541040e8c7db6, 0x99403e7365342e15, 0x457e9cee7c29cced, 0x8ece72dc073b1d67, 0x6e73cef17ad28d28, 0x7aed836ca317472, 0x89e1de9454263b54, 0x745329277aa0071b, 0xf623dfc73bc86b9b, 0xb8e3c1d8a9245882, 0x6ad0b3d317770bec, 0x5b406e8d502b}}

// C = 933177602672972392833143808100058748100491911694554386487433154761658932801917030685312352302083870852688835968069519091048283111836766101703759957146191882367397129269726925521881467635358356591977198680477382414690421049768*i + 9088894745865170214288643088620446862479558967886622582768682946704447519087179261631044546285104919696820250567182021319063155067584445633834024992188567423889559216759336548208016316396859149888322907914724065641454773776307
var curve_C = ExtensionFieldElement{a: fp751Element{0x4fb2358bbf723107, 0x3a791521ac79e240, 0x283e24ef7c4c922f, 0xc89baa1205e33cc, 0x3031be81cff6fee1, 0xaf7a494a2f6a95c4, 0x248d251eaac83a1d, 0xc122fca1e2550c88, 0xbc0451b11b6cfd3d, 0x9c0a114ab046222c, 0x43b957b32f21f6ea, 0x5b9c87fa61de}, b: fp751Element{0xacf142afaac15ec6, 0xfd1322a504a071d5, 0x56bb205e10f6c5c6, 0xe204d2849a97b9bd, 0x40b0122202fe7f2e, 0xecf72c6fafacf2cb, 0x45dfc681f869f60a, 0x11814c9aff4af66c, 0x9278b0c4eea54fe7, 0x9a633d5baf7f2e2e, 0x69a329e6f1a05112, 0x1d874ace23e4}}

// x(P) = 8172151271761071554796221948801462094972242987811852753144865524899433583596839357223411088919388342364651632180452081960511516040935428737829624206426287774255114241789158000915683252363913079335550843837650671094705509470594*i + 9326574858039944121604015439381720195556183422719505497448541073272720545047742235526963773359004021838961919129020087515274115525812121436661025030481584576474033630899768377131534320053412545346268645085054880212827284581557
var affine_xP = ExtensionFieldElement{a: fp751Element{0xe8d05f30aac47247, 0x576ec00c55441de7, 0xbf1a8ec5fe558518, 0xd77cb17f77515881, 0x8e9852837ee73ec4, 0x8159634ad4f44a6b, 0x2e4eb5533a798c5, 0x9be8c4354d5bc849, 0xf47dc61806496b84, 0x25d0e130295120e0, 0xdbef54095f8139e3, 0x5a724f20862c}, b: fp751Element{0x3ca30d7623602e30, 0xfb281eddf45f07b7, 0xd2bf62d5901a45bc, 0xc67c9baf86306dd2, 0x4e2bd93093f538ca, 0xcfd92075c25b9cbe, 0xceafe9a3095bcbab, 0x7d928ad380c85414, 0x37c5f38b2afdc095, 0x75325899a7b779f4, 0xf130568249f20fdd, 0x178f264767d1}}

// x([2]P) = 1476586462090705633631615225226507185986710728845281579274759750260315746890216330325246185232948298241128541272709769576682305216876843626191069809810990267291824247158062860010264352034514805065784938198193493333201179504845*i + 3623708673253635214546781153561465284135688791018117615357700171724097420944592557655719832228709144190233454198555848137097153934561706150196041331832421059972652530564323645509890008896574678228045006354394485640545367112224
var affine_xP2 = ExtensionFieldElement{a: fp751Element{0x2a77afa8576ce979, 0xab1360e69b0aeba0, 0xd79e3e3cbffad660, 0x5fd0175aa10f106b, 0x1800ebafce9fbdbc, 0x228fc9142bdd6166, 0x867cf907314e34c3, 0xa58d18c94c13c31c, 0x699a5bc78b11499f, 0xa29fc29a01f7ccf1, 0x6c69c0c5347eebce, 0x38ecee0cc57}, b: fp751Element{0x43607fd5f4837da0, 0x560bad4ce27f8f4a, 0x2164927f8495b4dd, 0x621103fdb831a997, 0xad740c4eea7db2db, 0x2cde0442205096cd, 0x2af51a70ede8324e, 0x41a4e680b9f3466, 0x5481f74660b8f476, 0xfcb2f3e656ff4d18, 0x42e3ce0837171acc, 0x44238c30530c}}

// x([3]P) = 9351941061182433396254169746041546943662317734130813745868897924918150043217746763025923323891372857734564353401396667570940585840576256269386471444236630417779544535291208627646172485976486155620044292287052393847140181703665*i + 9010417309438761934687053906541862978676948345305618417255296028956221117900864204687119686555681136336037659036201780543527957809743092793196559099050594959988453765829339642265399496041485088089691808244290286521100323250273
var affine_xP3 = ExtensionFieldElement{a: fp751Element{0x2096e3f23feca947, 0xf36f635aa4ad8634, 0xdae3b1c6983c5e9a, 0xe08df6c262cb74b4, 0xd2ca4edc37452d3d, 0xfb5f3fe42f500c79, 0x73740aa3abc2b21f, 0xd535fd869f914cca, 0x4a558466823fb67f, 0x3e50a7a0e3bfc715, 0xf43c6da9183a132f, 0x61aca1e1b8b9}, b: fp751Element{0x1e54ec26ea5077bd, 0x61380572d8769f9a, 0xc615170684f59818, 0x6309c3b93e84ef6e, 0x33c74b1318c3fcd0, 0xfe8d7956835afb14, 0x2d5a7b55423c1ecc, 0x869db67edfafea68, 0x1292632394f0a628, 0x10bba48225bfd141, 0x6466c28b408daba, 0x63cacfdb7c43}}

var one = ExtensionFieldElement{a: fp751Element{0x249ad, 0x0, 0x0, 0x0, 0x0, 0x8310000000000000, 0x5527b1e4375c6c66, 0x697797bf3f4f24d0, 0xc89db7b2ac5c4e2e, 0x4ca4b439d2076956, 0x10f7926c7512c7e9, 0x2d5b24bce5e2}, b: fp751Element{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}}

func TestOne(t *testing.T) {
	var tmp ExtensionFieldElement
	tmp.Mul(&one, &affine_xP)
	if !tmp.VartimeEq(&affine_xP) {
		t.Error("Not equal 1")
	}
}

func (P ProjectivePoint) Generate(rand *rand.Rand, size int) reflect.Value {
	f := ExtensionFieldElement{}
	x, _ := f.Generate(rand, size).Interface().(ExtensionFieldElement)
	z, _ := f.Generate(rand, size).Interface().(ExtensionFieldElement)
	return reflect.ValueOf(ProjectivePoint{
		x: x,
		z: z,
	})
}

func (curve ProjectiveCurveParameters) Generate(rand *rand.Rand, size int) reflect.Value {
	f := ExtensionFieldElement{}
	A, _ := f.Generate(rand, size).Interface().(ExtensionFieldElement)
	C, _ := f.Generate(rand, size).Interface().(ExtensionFieldElement)
	return reflect.ValueOf(ProjectiveCurveParameters{
		A: A,
		C: C,
	})
}

func TestProjectivePointVartimeEq(t *testing.T) {
	xP := ProjectivePoint{x: affine_xP, z: one}
	xQ := xP
	// Scale xQ, which results in the same projective point
	xQ.x.Mul(&xQ.x, &curve_A)
	xQ.z.Mul(&xQ.z, &curve_A)
	if !xQ.VartimeEq(&xP) {
		t.Error("Expected the scaled point to be equal to the original")
	}
}

func TestPointDoubleVersusSage(t *testing.T) {
	var curve = ProjectiveCurveParameters{A: curve_A, C: curve_C}
	var xP, xQ ProjectivePoint
	xP = ProjectivePoint{x: affine_xP, z: one}
	affine_xQ := xQ.Pow2k(&curve, &xP, 1).toAffine()

	if !affine_xQ.VartimeEq(&affine_xP2) {
		t.Error("\nExpected\n", affine_xP2, "\nfound\n", affine_xQ)
	}
}

func TestPointTripleVersusSage(t *testing.T) {
	var curve = ProjectiveCurveParameters{A: curve_A, C: curve_C}
	var xP, xQ ProjectivePoint
	xP = ProjectivePoint{x: affine_xP, z: one}
	affine_xQ := xQ.Pow3k(&curve, &xP, 1).toAffine()

	if !affine_xQ.VartimeEq(&affine_xP3) {
		t.Error("\nExpected\n", affine_xP3, "\nfound\n", affine_xQ)
	}
}

func TestPointTripleVersusAddDouble(t *testing.T) {
	tripleEqualsAddDouble := func(curve ProjectiveCurveParameters,
		P ProjectivePoint) bool {
		// XXX move this boilerplate to a function
		var Aplus2C, C4 ExtensionFieldElement
		Aplus2C.Add(&curve.C, &curve.C) // = 2*C
		C4.Add(&Aplus2C, &Aplus2C)      // = 4*C
		Aplus2C.Add(&Aplus2C, &curve.A) // = 2*C + A

		var P2, P3, P2plusP ProjectivePoint
		P2.Double(&P, &Aplus2C, &C4) // = x([2]P)
		P3.Triple(&P, &Aplus2C, &C4) // = x([3]P)
		P2plusP.Add(&P2, &P, &P)     // = x([2]P + P)

		return P3.VartimeEq(&P2plusP)
	}

	if err := quick.Check(tripleEqualsAddDouble, quickCheckConfig); err != nil {
		t.Error(err)
	}
}

func BenchmarkPointAddition(b *testing.B) {
	var xP = ProjectivePoint{x: curve_A, z: curve_C}
	var xP2, xP3 ProjectivePoint
	// This is an incorrect use of the API (wrong curve
	// parameters), but it doesn't affect the benchmark.
	xP2.Double(&xP, &curve_A, &curve_C)

	for n := 0; n < b.N; n++ {
		xP3.Add(&xP2, &xP, &xP)
	}
}

func BenchmarkPointDouble(b *testing.B) {
	var xP = ProjectivePoint{x: curve_A, z: curve_C}

	for n := 0; n < b.N; n++ {
		// This is an incorrect use of the API (wrong curve
		// parameters), but it doesn't affect the benchmark.
		xP.Double(&xP, &curve_A, &curve_C)
	}
}

func BenchmarkPointTriple(b *testing.B) {
	var xP = ProjectivePoint{x: curve_A, z: curve_C}

	for n := 0; n < b.N; n++ {
		// This is an incorrect use of the API (wrong curve
		// parameters), but it doesn't affect the benchmark.
		xP.Triple(&xP, &curve_A, &curve_C)
	}
}
