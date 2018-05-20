package visor

import "github.com/spo-next/spo/src/coin"

const (
	// Maximum supply of spo tokens
	MaxCoinSupply uint64 = 28e8 // 2800,000,000 million

	// Number of distribution addresses
	DistributionAddressesTotal uint64 = 280

	DistributionAddressInitialBalance uint64 = MaxCoinSupply / DistributionAddressesTotal

	// Initial number of unlocked addresses
	InitialUnlockedCount uint64 = 25

	// Number of addresses to unlock per unlock time interval
	UnlockAddressRate uint64 = 5

	// Unlock time interval, measured in seconds
	// Once the InitialUnlockedCount is exhausted,
	// UnlockAddressRate addresses will be unlocked per UnlockTimeInterval
	UnlockTimeInterval uint64 = 60 * 60 * 24 * 365 // 1 year
)

func init() {
	if MaxCoinSupply%DistributionAddressesTotal != 0 {
		panic("MaxCoinSupply should be perfectly divisible by DistributionAddressesTotal")
	}
}

// Returns a copy of the hardcoded distribution addresses array.
// Each address has 10,000,000 coins. There are 280 addresses.
func GetDistributionAddresses() []string {
	addrs := make([]string, len(distributionAddresses))
	for i := range distributionAddresses {
		addrs[i] = distributionAddresses[i]
	}
	return addrs
}

// Returns distribution addresses that are unlocked, i.e. they have spendable outputs
func GetUnlockedDistributionAddresses() []string {
	// The first InitialUnlockedCount (25) addresses are unlocked by default.
	// The first 25 addresses are use for SPO Network storage system,
	// Subsequent addresses can be unlocked at month by month. The rate will get fluctuation According to market
	// after the InitialUnlockedCount (25) addresses have no remaining balance,
	// except we have new features/projects which can add the values to SPO Network

	// The unlock timer will be enabled manually once the
	// InitialUnlockedCount (25) addresses are distributed.

	// NOTE: To have automatic unlocking, transaction verification would have
	// to be handled in visor rather than in coin.Transactions.Visor(), because
	// the coin package is agnostic to the state of the blockchain and cannot reference it.
	// Instead of automatic unlocking, we can hardcode the timestamp at which the first 30%
	// is distributed, then compute the unlocked addresses easily here.

	addrs := make([]string, InitialUnlockedCount)
	for i := range distributionAddresses[:InitialUnlockedCount] {
		addrs[i] = distributionAddresses[i]
	}
	return addrs
}

// Returns distribution addresses that are locked, i.e. they have unspendable outputs
func GetLockedDistributionAddresses() []string {
	// TODO -- once we reach 30% distribution, we can hardcode the
	// initial timestamp for releasing more coins
	addrs := make([]string, DistributionAddressesTotal-InitialUnlockedCount)
	for i := range distributionAddresses[InitialUnlockedCount:] {
		addrs[i] = distributionAddresses[InitialUnlockedCount+uint64(i)]
	}
	return addrs
}

// Returns true if the transaction spends locked outputs
func TransactionIsLocked(inUxs coin.UxArray) bool {
	lockedAddrs := GetLockedDistributionAddresses()
	lockedAddrsMap := make(map[string]struct{})
	for _, a := range lockedAddrs {
		lockedAddrsMap[a] = struct{}{}
	}

	for _, o := range inUxs {
		uxAddr := o.Body.Address.String()
		if _, ok := lockedAddrsMap[uxAddr]; ok {
			return true
		}
	}

	return false
}

var distributionAddresses = [DistributionAddressesTotal]string{
	"bRhJEVGunwNBs8Jtx2pogb831JqUhifokP",
	"2b1HJdnpAykBdZXyHRRGQ1tTRxFw1jgD1P3",
	"2ircJqBsANpor6LwssNZ9twgfuWKmcGoSk4",
	"2AFxdV1J1ZxjuXzeU2E1eEHyRiKFh2SVKHN",
	"uSuW26CuDNwC8HG4FbyxWpeh9pezpk9T1N",
	"2WenFpcN9T37kx1XnmTNnB2RUZsKE7hemQy",
	"bRCWuhGkyy7ScLcisbvp6zJXvbLD8a8HUa",
	"2Cvq9EqF3rYfRbiHJEwNjLYkn7Um86MAaP7",
	"2L7XTk3mNZXVgUzzG6MYxKyjssJkLsKyerw",
	"2QHJjJ5YPwjBJSu28HeoVFa1e2iTLiZmqUS",
	"2FkdR9nrxjUmt9LjdoFrqa5cL4kjXJ49G9X",
	"4hcXUcGFemfxNyF5LnVTQ7cHNDXFbVm2ZN",
	"2Zo3bqBVyyV53v6BpNdx4gVRSH4JXAXYkgT",
	"ELkipdznoXooV97aaKV9shmRMLqTmWT1tx",
	"2HjSKyjgJYybkSUitphbavQ2iueKowk3LGf",
	"ScneKQsBYfs7hgdj8v1xpkk9KU3Z5VPVXM",
	"2LmWA7wiLTbBJde1u7iux2JxLqMnKycxV1S",
	"2hqMK11rcRLoaA24z5bGjmY5AzWFZAbvuGR",
	"2YPb9N7pjXZhoG6F5tvB6mzGxhqdKz6DjFi",
	"BttA2aEDZdN7hR2NJMw8Vn5UKrgJxgWWGo",
	"27ebngLYti6KZHExM9PTXk7ZYibUEYDpPx6",
	"9HyLq9Dmd4KofoW6BPdmVEsa1nDoJekqNc",
	"2PJ9yVFapFUCHZChz5AGJC8Ums68APVBosc",
	"txmKDbEyj97NsrtGFKM4BxMsCg2bBC6iwy",
	"YaZRiVFJ2yGc5iVUvscAQh1hAh6ze7EN2s",
	"21LtZErKQsucga7JQy2chAGBhVM1NQarUTR",
	"eN139Ef2EXoAMdsJta778Cus27hPapeNuE",
	"TFvHQZyko4bAGtak58ScB1RkWQFxSRPfdM",
	"2TdgVEvki55wVCV15nT3qogLDLuQfpn26Fm",
	"2EkM32H9rrJBmEDCejJnUqUVMy3W71TfYFr",
	"hex7ZVQMmDPxZ6AEyG81Js9LksvLCn3XR3",
	"3GRCuqruvBCyFczgK3DMG6aLbtdMhJaB2z",
	"2STxoB8x93gg4sZ3S5F4C8DxmDxa28wtMnV",
	"2EzJ33mYUAcTESDB9utSUdUBs3bZ7289C4n",
	"ZMY5QESmpyqmMBy3aq9s8XSy6hE6zx1xdY",
	"o1X1SpFGiqPeGSp86JkozVeLV5AB75Vziv",
	"5N4wJPogpLoVbVm7tvgyWxwAZTW4NwRAVy",
	"APdA1Mq5qsJ1YJWd8YQsWXrNctYTn4BFP1",
	"R4vm6papR9ceNMrfHTG7qCERhYsTmDQ63x",
	"2SVKyTb5ehEuBAG7M2jakfYUZWGrQkdaoX7",
	"goYRfojX7U5ieLYKZpH2XT1HntgoUVxwDB",
	"2NgKwE1MpXZjQFGCSTxAcmxWuQ7mAckrHvJ",
	"Qd1xMSQXVVzYAFbxSLt4wF6paZp3eBGTRZ",
	"mqnZCTMkq5i7VdZreaQCbWpDr9vDHokmw5",
	"2JdFoCyjyvyosb8wRZmLU37Gqd36zcDuHNe",
	"phzYxNUvMXYScfPEVWKykA89JmvejsGUBN",
	"TC4qGd8CJaRsAP8GeDauucBpaGzcWEvroa",
	"22GAT5y2fuEZeHPkaUwthp23gjfsnqrc7nU",
	"RcYwRHG5nkET4ARAVEfKSjZSGMEBbTAqQ8",
	"2Djm1tiP7ddNHPefFmiPU1bTBtJ7eXTVVMm",
	"Pcz1PWtqgUBsFrzPkgkk8k7vk32hdkWzvf",
	"2MT6yRa1Z8oWULBcYkbNXN2cAZB8xxmW357",
	"qedFwAPhSN7NMRtP2fKb1GFNAtfTLjF24r",
	"znLTAU4c9VprwquKRKzSFTBNL3boMM8BXW",
	"wtsoA8DPxvEukbZLDAtdxhqkHaLyspoSUG",
	"2eehZgS6UQt1i7CJoMzKcEuM6jsFi64T7ua",
	"J1ZP4XwL9wUSgGmPzr3QFj92m7FgirsDp2",
	"nrffPNVMDb2tWwMnnhgoRDEHhJLaQBAsi8",
	"2FUnHArqyQDMsKcrRznWRWX47tmvDi8BbPv",
	"2BToDyTnqYH5aZxzJwvCYoLBiE2mRdsApc",
	"T5pcyi7r6gQWM1Z8bg6zH3z1PtKccmiZeq",
	"2T1Gbob6z6Ck62njTf7dZdapeLGWKoi67TK",
	"VCzjQFmMoxjiqNXaAGYRsctZGYbEYVuJmW",
	"2dVMrisaPHau7nikUUxzkCLt3KaBhxB77vn",
	"eg7NuKXWTJVyZ4Rcr7WQxK2qzUoHzBArTR",
	"GzFhfRdz2gMhNvnCueJ5h4WLPjE3TvKC58",
	"hsBTLx1SV1aj72KM4PTrFBSKsDEbt2XYPw",
	"4feEcntFEX2BUD3EVUbzQTBhVrztgwq8oq",
	"2NWqPi5vVHc8hY4D4iK2yUYF331TzNuVtpT",
	"2PawH4bvZZcrQdc7ERp1eAc1B9CpcXaMmRi",
	"FiZTj47GUCpUQ7ayCqTDxNpR88ABscBArJ",
	"LCJw9BrxsjdSSzQUbs7sCywS1TTDhGWtRu",
	"2ACXcoHE6FVb2dd9nh4igh4rWjQtVQDxnHX",
	"8gXDWYCjX7uLBGQbADa1VY552MqD8eeAbQ",
	"E9gBBn11W22F8vPVaQMnMuT5DPv2L1upkv",
	"2YFU2HWB9kcdz7HteRYK2G9ZV3M6Tg3x47e",
	"2HZVqDvVxPjQLZ2rHAEsa3f9UfDCsBftPHu",
	"tiyAoZ6WfBgacKNjtWhJwVfLJcy1CrBxjy",
	"adnpVEsprbtdit3QTxJ6bXtaNmxQpBquDp",
	"PCSHdzqbHKy3y7WJkqfcMkmZBMqtYEZiu7",
	"n1aCioquWEdZ4ajHmX2jyGtwP95f8F5p91",
	"2i2DqDfADp43Ue1gPvf94C9ZpkdW6g5mWid",
	"kRMDBxu2soduQe8itaWoNRj5yX28NMSqFb",
	"2moG8hfWUcgcXTDr5pNSGEqzEzKovJHbfV1",
	"Ech7gHemCpPYNDSxJA7nuTL3QecbX9LXDR",
	"HGQpFLkAvKGTdQqRht1NVU6LfGkg1end8o",
	"PiVEfGwixeECxD3PHuVAYxka4d4b3aJV1w",
	"Ldvitx3or7HuvE5np6a3xJNT3Vq3MQHnRg",
	"2PVD39YQXhCP7S4jizAdQQxfXjubUdKvTpL",
	"2A18amgqgAjVBgC1GugGauP4Vy9w5WSEaxw",
	"txuAmS9atjUwTg38HZHifAN9pVZ92WNegu",
	"KQuVUrfCyRY5rwAtHdzeWn19YVWCPE3ZBW",
	"2P3xosriAwNi7fgKgNJc3EaGpVUNhGPPFRr",
	"ghKfz92v7RanF2xVrE5hD3C6JwemZQgNn",
	"2MDLetXUmy3VQfHFhMA3zMQbZ7WNgNC4UYJ",
	"5yjxgEPzu3tbRaum4sKTJy7gi1jgDC23dh",
	"2Lktp8Q9jz3oebwDKbhymE3GAVvoELhVoms",
	"wESxGi34RKJMozknduJiUAmTpwKbnqPMtc",
	"2S6BMLzD6Fm5SLzKS6E2JvWdjRH2RVmiiE7",
	"UXYmB3W5MPkmHjCFpZVPnRjaB7FiFEJtH5",
	"2DihydNzykzYq3D6AWAjafUgA3FATm28xCd",
	"2gchvg3qenoHFqHh99N7CjErSBqFAft7Bhk",
	"24SK8vz4C7P6MC4rkTrVCyMwhsQG6yRGy54",
	"2SaDXt96YfA1EtcECRWv184dhAVRdH2ueWU",
	"xAjqp8TZvbLF3vtgTq8nyyM8h4Y8gSnHG6",
	"nMW7XZUvCQSnauyK3ByujPAU3nfCUdy3Pb",
	"2VR7udxTUDw6wXrkq172uUpLifnX9KaTjTG",
	"xmiDRywB1LeApPgXg2enSsWwioeqjbFuWk",
	"2WuUCE5RwiS4UjF6Vt7ZWC9MfcyuXH4zDpb",
	"2KhqVuvtkzA7DU922aymSW9VF7X4iuNCxwN",
	"xmcknTven7YckxpkodtBZkwfa7NVFFsrYF",
	"vfGNYtUqdgqfFekLuHUqFp6BdKh6JyLder",
	"fAYG5ESqjvLMMuspjK8zXutUBP3emCt8zn",
	"D3t6D1GAVmCEJ7a1PgKkdoGsw6vyCaPt4K",
	"2fSQ4yAhP9iu9NVe2qvkNGQHWaB6V5D12MW",
	"UrBr2NKRfLZkhozJ66nEz4WcDWGEDUWKgV",
	"MeqZU79PdUwHkzfQS6NqhMUwhP3qNJcE1o",
	"fPRAQyp92oAZW34ocqKqtuVEBAHBfgS6i4",
	"2aCae8EQN8hT76P9CVLk6mP8moivJHc5EkN",
	"24V9j23tkVkN8yyFtsG3D8aLZjn6C4ug4ed",
	"TY9E49vCt97LQPSAbChApRmpdeDGaeTPtG",
	"yihg6wMUsVFrHpMn3hvF3JdkvpgihmaC8K",
	"ujwaugBsr3rSZKXrRRiqRCXhxJo8rP32v1",
	"aBmokr2v6BrbWxDMrVQM3Tci82LKtPtfVx",
	"ZZWa6Tqxv5efutK28yMJyJnJToS8nQ1QNW",
	"2Fivgbf2SThrvsDfTow7uPANtuqQ1TRtdbz",
	"oZAsNiDuJb4w2kHwBbntvAj2nRwwiSTFZ8",
	"uzDiiBq3DeWcNf9AAc3vhcxEkyxNNQfZGa",
	"SKwEQXCEzEGjpFhDtpJgB3htRcns2DjrZQ",
	"2TrNDnQB6QxsUiJAFJzDHihQox6rvw1Fhh7",
	"oZAxtUzXCLXWsb7bHxXAtf9fQa9VVg4R17",
	"2kS8UFJPr4k2rDCHrKnqEyY5KCHATn3WAWF",
	"L1az6BZf5m2LiAgF411ArjE68P7bE74dF9",
	"ChvMaNFLadB7FTdMHDTDbZL5xRBsrUqwdY",
	"2gWr7224GinDQwHPkCjFjuUNmDBnQqf3qNX",
	"sZS79d1GuiJXKfWixP1cpVf1VPXZKC6bry",
	"ZTDLymFo1zadXpF9CbvFUZGPj69w888wMZ",
	"2UWrs4s2zAp3pa4s1ueCMwdFeazi9zzev5J",
	"2XrUn6BsRundm82mEMF9d4BZX6gjmnMkpqt",
	"2Zd4dKceiryoTf4kahuhreZYFqLhGwGWQQt",
	"6MawZSRPpSSokKTMrXbeBapyw5fH5AGkDP",
	"QyHvNhLaaUTPFoY7U9gBeiUH24F9SYCLwW",
	"fB5HMXDw8nTC2Th2uVCvLw3yacmnbXvnWX",
	"znhYPuevKtmuKYUD8E1PYwZomjKJNKUx29",
	"2aRFMGHvGAkeoSYDLvKwD2uCJE5AizKM1p8",
	"8HSNRJHy5P2PpZYWHfbfTjkz7mcGNCPp5f",
	"2Mv9x4AmgvEMyVgvSfooFuUXdoJcNEAgsXS",
	"2HdgxEJDxARgSQFZZAHkKHjkkQWa8PW5Qk3",
	"UhnGXHLcSacDnE13n1qwUUNnJAZT1DQgvD",
	"xg16s3oWVMfdDf3176PZdHAt6kYV6gtLMU",
	"2gWsKHVr8a8yNbVWms9KvdZgZuL8EjYyAir",
	"wVpDHATU4XQ3dAUVxsB2kT64T17e2r49Ji",
	"2Dyn136RN8HtDPMEDiAUgaiPEU9rdJQP1J1",
	"PbJW7wn1VcJisWNXar3w4HekXK7sr4SnvX",
	"2U5u1mu3YSDNnKSwbUuVMSQCVoTx2RdjBPi",
	"AnrEKKbRNbeYCf3F5dgGQSYaPUQioSDrpK",
	"rLK3gumhz7mLRDSTxwCmSEMkCtXHk4PfJe",
	"rqzTPARP6WetzniiNSP3AeCGYNMZU1JfRG",
	"Tke4XKKn4SLGn9c4C3g9UNJPSunSgKheTM",
	"2bjZ1xGaBQJE18jRPRmxaGA2DXXgM8NsUJd",
	"o62jYbfb3BCCXr1SD2myhF2T5TNL3t7dRz",
	"6Gg3XtenQAenenJci43DWLe5LaehAVrvTL",
	"Ye8dZvqWQE6eYRbpbj9wJZWVtQwbRukNJi",
	"x3XZvQ6oFxnjwsoQy59MwqEwSNWDWaBfSi",
	"2RPuXJzQmsMfg2YMepxEbYeX8vmJQCbD6Rg",
	"2m3oWbuqntMRGvTpjKjvv1Sjk21rfk57naS",
	"2JATanSkpcSDBMW56xvb4D6Ujm4eHBmrqsg",
	"2FdoA3gvzNDz8L64WAKzWPm818cA98oLnC3",
	"2mqet3fPE7h1kGCui9nWcQxLw37xSv3KF3Z",
	"2CXeEBzHJSnyaYUnR5z3ca1eN2egM49M8z4",
	"pcaFwMxNvw6fZQR6uTaGiSrgho8VxE4nu6",
	"hgHbFcYBYTyzVwfGCVwZKVB4ZnLVdNVnxH",
	"26Esxno2M9woXc1PjKUHgJp7DJo6Azrd4Bw",
	"sfCJvxJ7Tz81iTcJGY9g8HkcDWkbPMf7M5",
	"25gwdyUXFsCtbgt5LKY26xhieH24kmjM81L",
	"2erxRBLUuxWHmJu6ZpfMA6HM797pDKDi8AJ",
	"b5EJUivJHg7ctu1EyEanoRR79Cds5TDSU8",
	"2PeJEtunhggrwD2reS97RtBXuPnZMB5FXHr",
	"qP5Fx2sqq4SVQsUwQ9HB6tr46yyFzaLhwT",
	"2FSuboSpWgj8Vsb9J5zbTqhi6se3qDepGw8",
	"8FDf5qSUSL8PiycSaoyfomFJyvShPVf6hB",
	"2gDxB4CjpV6fURvdnaJ4rCAHLkb45vnqYbw",
	"1JuMh1mhLwtf7NwTwtLu8vDdWzySY98d4F",
	"2RcnnDogtRaTqDBA4VnhTumHe191qbbupfu",
	"YYeqcyhPHQ5BSESMKK4HwdszySV1QLjffv",
	"2bWpNz9U8g8gSFHdmWzLo5NwkLSecX3uNrZ",
	"2NUU6hKaT46Assf1ukuJ5mhqb2NEYbpyZd",
	"58yJNafupMdEKyGgEppp1QVYGD8GF36i64",
	"xvvQqxKxFVyk2HuPLaRJnmKw33EX1ohF62",
	"JYWcSLDEYkwLnUzNPCAbbheCDgBX8UBoWz",
	"qqjgkoip9uV7BCkrSJ6j4FeafzFF2c6U9b",
	"EiL7qGeMFUL8gxDfmKTFBBybevAsabPjQs",
	"6YX7BXbUFvz6TFFEeZiRNNFematw2Ndy7H",
	"PoaSVBs78jvVYZhZws51qsPuSMoVzj7P2h",
	"2LSofXm3o3gjubaWr9uLsZiVEHznid9nj2s",
	"uW2h5AWA7fWXhWCSFe38Qt8ZfaAPxTCFdc",
	"LKMcUWtkRAZ7RXpjcpw9r2m8wipf6bdzPx",
	"2JNiM2pWip6azRUMrMEYWaU5mwjWAyPbB1B",
	"29ssxDkFEYoim8fs5wMFBau2xwk5DojKDpb",
	"2RSdPLfAmqBNg7WQtTAKzDzavKu9kTq7Kzb",
	"nNBmjY1Mgvm9NzfR3f2RVNBGiHtCDh6CHy",
	"2PvqrSzfK81rLwAMnzYoA7qa43ekBrch5c9",
	"2UM2EYCTqxBNJGcccbDKmyGhshuggzP11dY",
	"2S8FZ4geCbBfyxrMrkc71KVixFhvVjG4Frd",
	"GnuLHQYrS6H4pD87tuFZ9bcGxbfeJdcLEx",
	"2R1ogMmJ3qdm9asPTuPuzo4zQpttfHA1aaX",
	"LVK3BWnGsy7Rs15uR5WJzawRyq3JkuxHmg",
	"QAE88CYhRUuvSijqKNxodBWNBo61k7Yyra",
	"uQNrnbfwPVeVEZRKizW9pXqpnr8g3TkaXU",
	"2Me8ruE32ok73gRDbUKqHYGLTT3ua3z6Mmp",
	"2cWx6ifxM7NWcekyswKwpQz91JaJ2rJAxPJ",
	"4Szp47p8c2Rv4M9FJYqmz5V32TNn45YR1P",
	"XsHSFqWP9554b7GLVE7R7RGyaWr3WqX9g8",
	"ppLuAmsufVpGbNQfUzMJ2Hc51JkCX2tMK4",
	"2JsFtN6Ysvb5mbmJkwzcLGXYzFHNQQtGpbu",
	"2RJRox6aEEpZcorAzf43RBAjh3iatrNro75",
	"txDFN3URsxgimLYR8ALFKRuiwiitnQpPEy",
	"NvCXCndppZEYcx5GeEvT2cWurAsUCPZPDq",
	"2Vvi7WSDGXVqLUcKLNzErfJ3koSywoWwzRy",
	"rESHfZNaodMfdTdJdhzXnWapTdKR8FcPYR",
	"2JKqDUcY6ptde4j2tRqJFRFvFSk3yPzEG8N",
	"2ECZeVS9K2zebayY6f7K2vtgdmSTNrc2eFc",
	"b5bgMPaoCRbLBNyaVVsk8ZC1RWeLGr5Gq9",
	"7Z4WvHkGNAeQ3w2XEzFHgHAwTLTvDSvQ7Q",
	"REsQgCEFkaQnxhesUwaadcFn36sKJPzsqF",
	"dm31zSQ52nMGgxneSPUw2VfNe8EFkeLHZZ",
	"YPDhXjEi2Fomw6gHWPTAdQhDRByTwGxV8x",
	"K73AHZnivQbugR4BodRFL14dghcQ4WnF8n",
	"KYABXLbraxs2njmptW3XJswAWDwCAQaH7N",
	"2j8cwkGeaLJJbXycGAECa5NUN6wjEjnEh4U",
	"8dV6ncpKLi6RoUN8HXZo5YL4QWftJc1Q9L",
	"W1qrPh4dyQA4ZoHtqQTRk3mPY1aFPvUbHG",
	"xAhKZud7ByNHApgSxiXeR3weG1iR4R6LFG",
	"Fhw5ELwa1HZpbAo3UsqaCCRqd6qApJYvpN",
	"2DTX8A49bW8MPo5576t6Hg9Nh9RGRKWLaXo",
	"RQS9j2Mrf7uAz9hfdJc46N3Q8N2WNsHsbU",
	"KzLooZtBMt1FwhvQ8kBPDUt98Z2CKwrHUV",
	"2VWrxk1pV13LSuXczYNwbBc2dvZQSdg4S2b",
	"21rJ2poqarhUNaUoTmLThiEmSxFWx22w9Ds",
	"qnim7Sp76BCt7UV5dQDycVChauPwUFmf2M",
	"2PotPMW64mg8kx1RtofCwqoZ8H4S3QDDv3T",
	"2ZNxiKGfmWWUhv8BDxefqgxnbf8ap9mXyiJ",
	"2hSzpZPrdqdkZeZBumGcBvCzRyqDY927WAw",
	"21K6eJuS6V29LVAqo3WBcmjQT1zzHC8KHMt",
	"s2brAnug39AM9CKDaUn1dtt2X3aczaFbwg",
	"GmAxTjLGuqzaf4oTqZzQvmuhr6dww4c2ct",
	"2AZNEk9LEW5kh2xDcKtjXfGQqc1LDxrTLcL",
	"27a5s7BBmopXoNeHzjv5n3vtK55Mz4TQDwY",
	"jbay7JcQkMTuid2pxiLtvCq8eh1h31rmWb",
	"2WcTVoS3jGBAQkYHaQs9WJ295gyurM2dY4E",
	"2GJBxPuvG7JvvGrLfwBNqgSTinYuWZtbXRx",
	"2dfjoBbpt6eWVFwtFxNfXL4RhWg391feWSG",
	"25pxfk5nyPHMFFsSMp5AESDW9gLPVeJTes6",
	"2PsfV3WEwEi3436dXe22174f1YzCY4gaAAU",
	"PvsdXzXmYjb7NiQkUKb3JMkcfMPodhw5z7",
	"Rcixy3jNJ3pSNS8YsM1ktfoGHjwnZ78MrQ",
	"DhVndspgugvY5QL8NiJLhNrxX1XukNNyqL",
	"c9DH3m5uJV1NMPSS5EAyBSgAMW1YvciNLG",
	"Q9p514badWQChhFc9oyn4VZ3HrfVj4rtxm",
	"PW9utoNCh1PDpP2QZWKy95Q3WSHW9HBn51",
	"Zr6NKwLNyYs93rsEsPwSvijZ1abHn5rfhQ",
	"2mdxeXdQigw8sWKAJeVXnsxhxUnMNr5uoHe",
	"256RG47K9RpbFoFKW7NJM8iHDbTLWS24Gf1",
	"CzGPN5BEXSi8QEApXGYKb7Q5nxRstkc3xo",
	"cu5SJKdMgTZ2HqZqGUCqePjsKuqZXwrxod",
	"eLaHmMkJj5mqRB82ocwxR8n6FQYLbKqKkg",
	"qyeDQXZXeNcc4dQup6qrDW2oHWLTpEv7Tn",
	"6AGXcXBa8sztDEmvD3X6hRWVVdcBgZNgng",
	"zgevaKxfdGBYzdLXsZTwG83UNiXpHpfp3m",
	"QVxjqMRVVzMz6JRfqZUWRJ4yu5LTQZeaWj",
	"LdPdEgtRcSDTFuVgDRguT8cD9XU4jgug5r",
	"25kan7jcfLQtXMPV7Qx8Rwu558Hw8QVQHVn",
	"7EobzC9YPKJyovinAb28sM8cMVQCY9ji2c",
	"2QK5zXUoDqyAd6zX1WcvwDrfGwNnq2anQnJ",
	"2VbhDRiNW4mJxU2wfyj7GR2ywRmC7xEYHUs",
	"ppLjPCDXonLHhGZXna51GWtLXjvnfG5PeP",
	"2FjHYQ9hTvxN8oXzvStqjZqA58s4FUQM6Gf",
	"h1e51zqKSTYzFRagmAF1R81Wr4kVHB2u57",
	"bGqjqCJdKqE3PrWGNnFrAQtap11XyRF3Jt",
	"2fjEV5tJEsg2Ppxa9A8nY8mWi3ZHCZk893c",
}
