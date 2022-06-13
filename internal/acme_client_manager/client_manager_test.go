//nolint:golint
package acme_client_manager

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"math/big"
	"testing"

	zc "github.com/rekby/zapcontext"
	"go.uber.org/zap"

	"golang.org/x/crypto/acme"

	"github.com/rekby/lets-proxy2/internal/cache"

	"github.com/gojuno/minimock/v3"
	"github.com/rekby/lets-proxy2/internal/th"
)

//go:generate minimock -i github.com/rekby/lets-proxy2/internal/cache.Bytes -o ./cache_bytes_mock_test.go
func TestClientManagerCreateNew(t *testing.T) {
	e, ctx, flush := th.NewEnv(t)
	defer flush()

	mc := minimock.NewController(e)
	defer mc.Finish()

	c := NewBytesMock(mc)

	var err error

	//register account
	manager := New(ctx, c)
	manager.httpClient = th.GetHttpClient()
	manager.DirectoryURL = th.Pebble(e).HTTPSDirectoryURL

	c.PutMock.Return(nil)
	c.GetMock.Return(nil, cache.ErrCacheMiss)
	manager.DirectoryURL = th.Pebble(e).HTTPSDirectoryURL

	// create first client
	client, clientDisableFunc, err := manager.GetClient(ctx)
	e.CmpNoError(err)
	e.NotNil(client)

	client2, _, err := manager.GetClient(ctx)
	e.CmpNoError(err)
	e.True(client2 == client)

	// create client if all prev clients disabled
	clientDisableFunc()
	client3, _, err := manager.GetClient(ctx)
	e.CmpNoError(err)
	e.True(client3 != client)

	_ = manager.Close()
}

func TestClientManagerGetFromCache(t *testing.T) {
	e, ctx, flush := th.NewEnv(t)
	defer flush()
	ctx = zc.WithLogger(ctx, zap.NewNop().WithOptions(zap.Development()))

	mc := minimock.NewController(e)
	defer mc.Finish()

	c := NewBytesMock(mc)

	var err error

	manager := New(ctx, c)
	defer func() { _ = manager.Close() }()

	state := acmeManagerState{
		Accounts: []acmeAccountState{
			{
				AcmeAccount: &acme.Account{},
				PrivateKey: &rsa.PrivateKey{
					D: big.NewInt(123),
				},
			},
			{
				AcmeAccount: &acme.Account{},
				PrivateKey: &rsa.PrivateKey{
					D: big.NewInt(222),
				},
			},
		},
	}
	stateBytes, _ := json.Marshal(state)

	c.GetMock.Return(stateBytes, nil)

	client, _, err := manager.GetClient(ctx)
	e.CmpNoError(err)
	e.NotNil(client)
	e.CmpDeeply(client.Key, state.Accounts[0].PrivateKey)

	client2, client2DisableFunc, err := manager.GetClient(ctx)
	e.CmpNoError(err)
	e.True(client2 != client) // get second client

	client3, _, err := manager.GetClient(ctx)
	e.CmpNoError(err)
	e.True(client3 == client) // get first client, cycle

	client2DisableFunc()

	client3, _, err = manager.GetClient(ctx)
	e.CmpNoError(err)
	e.True(client3 == client) // second client disable, got first

	ctxCancelled, ctxCancelledCancel := context.WithCancel(ctx)
	ctxCancelledCancel()

	client4, _, err := manager.GetClient(ctxCancelled)
	e.CmpError(err)
	e.Nil(client4)
}

func TestClientManager_nextEnabledClientIndex(t *testing.T) {
	table := []struct {
		name             string
		accountsEnabled  []bool
		lastAccountIndex int
		resIndex         int
		resOk            bool
	}{
		{
			"Empty",
			nil,
			0,
			0,
			false,
		},
		{
			"OneEnabled",
			[]bool{true},
			0,
			0,
			true,
		},
		{
			"OneDisabled",
			[]bool{false},
			0,
			0,
			false,
		},
		{
			"TwoEnabledLastFirst",
			[]bool{true, true},
			0,
			1,
			true,
		},
		{
			"TwoEnabledLastSecond",
			[]bool{true, true},
			1,
			0,
			true,
		},
		{
			"TwoDisabledLastFirst",
			[]bool{false, false},
			0,
			0,
			false,
		},
		{
			"TwoDisabledLastSecond",
			[]bool{false, false},
			1,
			0,
			false,
		},
		{
			"DisabledEnabledFirst",
			[]bool{false, true},
			0,
			1,
			true,
		},
		{
			"EnabledDisabledSecond",
			[]bool{true, false},
			1,
			0,
			true,
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			e, _, flush := th.NewEnv(t)
			defer flush()

			m := AcmeManager{
				lastAccountIndex: test.lastAccountIndex,
			}

			for _, enabled := range test.accountsEnabled {
				m.accounts = append(m.accounts, clientAccount{enabled: enabled})
			}

			resIndex, resOk := m.nextEnabledClientIndex()
			e.Cmp(resIndex, test.resIndex)
			e.Cmp(resOk, test.resOk)
		})
	}
}

func TestClientManagerStateMigration(t *testing.T) {
	t.Run("FromVersion0", func(t *testing.T) {
		e, _, flush := th.NewEnv(t)
		defer flush()

		initialState := `{"PrivateKey":{"N":24765388055830593543913638614972357269643277023023551760574908381166575195664233540375003284842329984712435378951124200758916192940220675139985244986481184467176347315627474926703581063619199463482962612359663622026588917057257145909601319955399123550101448497320956972704611297541803861338402199184188287717597438874232925988270352397747350594809299325379334836529433414400861008826247004652135118227836570919945934737897184246008549742424511866623479010090267698061008506880566222992529934956196807614325622913886973860017827185835907076338726896147347335231597743074530442788894304599011771252372158894921500658841,"E":65537,"D":179117047445926748856601075781572201135403105252195912759395556291453021083431446330130331828055364340047520784027844899213059423740247494031661597625647823938257604522749334196827706854990319143246170533568527043358761412410392406749638000806554840208555267829319828571066508308815097277482988882818945761600188258009545759600458645762417807389718958086576424290725153851197170643671061809097074314385705628355964886673229085614667965931410510440706671327119787337495299059814524268506060542887248989836486629377936966940896622020698260241255184996595631192278639788160773584801277528416318321001164730636495622473,"Primes":[170240283176194179466016930790115076629185096952978648934338872566459814202728254055430647846801527733066683640097070577744691771303100786417278209859008005833233224181812896410107081093341807361831899160417102679896129803828108888686970220007896733517406192240896554203960573075774658133241012083757109546551,145473137108207659151891864991885638644225291464417118520718153762863278011133987341862784005727047900877414013615694714828215601402316939377945275409636268731088754253469619862841229813906398887696626222790852187704942282070030505307577226515214564903294471093776043801977348865388476402984874761113729487791],"Precomputed":{"Dp":102361850540017209149608056131883893215277903024689513403215856880142750484042748055423655477837877868215294192924075914164629506080076744275131073099227451330765650428132489981791654142915105374068813316376952686329659445040976225149133290959781172177348801615038369393796944361566385071005620064582337592373,"Dq":39220056755190520309974172796155261138056619541400828038246624942185807393971747598202138190597543515275386853328283411146371385037117017896286297377250457485231353264637604915643707975371349954174156411347659601516069169810875825843105146944433314727197003368416755694296928103247768276917462687552351897703,"Qinv":65224836018493488025550562263913342772761002765710873750661002958212765185683954846872486718530224759527254500982976806738524916536298648798971146903631592569092848720420517387339362811827466959138620567561472378959342036278564222532139577521100078498611736513180404113216987840961790328448144484052030954491,"CRTValues":[]}},"AcmeAccount":{"URI":"https://acme-v02.api.letsencrypt.org/acme/acct/485823100","Contact":null,"Status":"valid","OrdersURL":"","AgreedTerms":"","CurrentTerms":"","Authz":"","Authorizations":"","Certificates":"","ExternalAccountBinding":null}}`

		var state acmeManagerState
		migrated, err := state.Load([]byte(initialState))
		e.CmpNoError(err)
		e.True(migrated)
		e.Len(state.Accounts, 1)
		e.Cmp(state.Accounts[0].PrivateKey.E, 65537)
		e.Cmp(state.Accounts[0].AcmeAccount.URI, "https://acme-v02.api.letsencrypt.org/acme/acct/485823100")
	})
}
