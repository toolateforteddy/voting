package vote

import (
	"testing"
)

func Test_firstPastThePost_Winner(t *testing.T) {
	type fields struct {
		votes map[Candidate]uint64
	}
	tests := []struct {
		name    string
		fields  fields
		want    Candidate
		wantErr bool
	}{
		{
			"Simple",
			fields{
				map[Candidate]uint64{
					"Sally": 5,
					"Jim":   3,
					"Ted":   2,
				}},
			Candidate("Sally"),
			false,
		},
		{
			"Simple",
			fields{
				map[Candidate]uint64{
					"Jim":   3,
					"Ted":   2,
					"Sally": 5,
				}},
			Candidate("Sally"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &firstPastThePost{
				votes: tt.fields.votes,
			}
			got, err := f.Winner()
			if (err != nil) != tt.wantErr {
				t.Errorf("firstPastThePost.Winner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("firstPastThePost.Winner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_firstPastThePost_Vote(t *testing.T) {
	type fields struct {
		votes map[Candidate]uint64
	}
	type args struct {
		b Ballot
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Simple",
			fields{
				map[Candidate]uint64{
					"Sally": 1,
					"Bob":   2,
				},
			},
			args{
				b: &ballot{nil},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &firstPastThePost{
				votes: tt.fields.votes,
			}
			if err := f.Vote(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("firstPastThePost.Vote() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

const (
	ally    = "Ally"
	bob     = "Bob"
	charlie = "Charlie"
)

func TestFirstPastThePost(t *testing.T) {
	ballots := []Ballot{
		&ballot{
			[]Candidate{ally, bob, charlie},
		},
		&ballot{
			[]Candidate{ally, charlie},
		},
		&ballot{
			[]Candidate{bob, charlie},
		},
	}

	fptp := firstPastThePost{map[Candidate]uint64{}}

	for _, bal := range ballots {
		err := fptp.Vote(bal)
		if err != nil {
			t.Fatalf("error voting.\nballot: %#v\nstate: %#v\nerror: %q",
				bal, fptp, err)
		}
	}

	winner, err := fptp.Winner()
	if err != nil {
		t.Fatalf("error getting winner. state: %#v\nerror: %q", fptp, err)
	}
	if winner != ally {
		t.Fatalf("winner not Ally. Winner was %q", winner)
	}

}

func TestApproval(t *testing.T) {
	ballots := []Ballot{
		&ballot{
			[]Candidate{ally, bob, charlie},
		},
		&ballot{
			[]Candidate{ally, charlie},
		},
		&ballot{
			[]Candidate{bob, charlie},
		},
	}

	approval := approval{map[Candidate]uint64{}}

	for _, bal := range ballots {
		err := approval.Vote(bal)
		if err != nil {
			t.Fatalf("error voting.\nballot: %#v\nstate: %#v\nerror: %q",
				bal, approval, err)
		}
	}

	winner, err := approval.Winner()
	if err != nil {
		t.Fatalf("error getting winner. state: %#v\nerror: %q", approval, err)
	}
	if winner != charlie {
		t.Fatalf("winner not Ally. Winner was %q", winner)
	}

}

func TestSingleTransferableVote(t *testing.T) {
	ballots := []Ballot{
		&ballot{
			[]Candidate{ally, bob, charlie},
		},
		&ballot{
			[]Candidate{ally, charlie},
		},
		&ballot{
			[]Candidate{bob, charlie},
		},
		&ballot{
			[]Candidate{charlie, bob},
		},
		&ballot{
			[]Candidate{charlie, bob},
		},
	}

	stv := singleTransferableVote{map[Candidate][]Ballot{}}

	for _, bal := range ballots {
		err := stv.Vote(bal)
		if err != nil {
			t.Fatalf("error voting.\nballot: %#v\nstate: %#v\nerror: %q",
				bal, stv, err)
		}
	}

	winner, err := stv.Winner()
	if err != nil {
		t.Fatalf("error getting winner. state: %#v\nerror: %q", stv, err)
	}
	if winner != charlie {
		t.Fatalf("winner not Ally. Winner was %q", winner)
	}

}
