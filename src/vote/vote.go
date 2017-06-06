package vote

import (
	"errors"
	"sort"
)

type Voter interface {
	Vote(Ballot) error

	Winner() (Candidate, error)
}

type Ballot interface {
	Vote(...Candidate) error

	NextChoice() (Candidate, error)
}

type Candidate string

type ballot struct {
	candidates []Candidate
}

func (b *ballot) Vote(cands ...Candidate) error {
	b.candidates = append(b.candidates, cands...)
	return nil
}
func (b *ballot) NextChoice() (Candidate, error) {
	if len(b.candidates) == 0 {
		return "", errors.New("no more candidates voted for")
	}

	next := b.candidates[0]
	b.candidates = b.candidates[1:]
	return next, nil
}

type firstPastThePost struct {
	votes map[Candidate]uint64
}

func (f *firstPastThePost) Vote(b Ballot) error {
	cand, err := b.NextChoice()
	if err != nil {
		return err
	}
	recordVote(f.votes, cand)
	return nil
}

func (f *firstPastThePost) Winner() (Candidate, error) {
	winners, err := getMostVotes(f.votes)
	if err != nil {
		return "", err
	}
	if len(winners) == 1 {
		return winners[0], nil
	}
	// TODO(teddy) Handle ties...
	return "", errors.New("tie")
}

type approval struct {
	votes map[Candidate]uint64
}

func (f *approval) Vote(b Ballot) error {
	for {
		cand, err := b.NextChoice()
		if err != nil {
			// Improve error handling.
			return nil
		}
		err = recordVote(f.votes, cand)
		if err != nil {
			return err
		}
	}
	return nil
}

type candidateWithVoteCount struct {
	cand  Candidate
	votes uint64
}

type byVotes []candidateWithVoteCount

func (a byVotes) Len() int           { return len(a) }
func (a byVotes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byVotes) Less(i, j int) bool { return a[i].votes > a[j].votes }

func (f *approval) Winner() (Candidate, error) {

	sortedCand := []candidateWithVoteCount{}

	for k, v := range f.votes {
		sortedCand = append(sortedCand, candidateWithVoteCount{k, v})
	}

	sort.Sort(byVotes(sortedCand))

	if len(sortedCand) == 0 {
		return "", errors.New("no candidates")
	}

	return sortedCand[0].cand, nil
}

func recordVote(votes map[Candidate]uint64, cand Candidate) error {
	val := votes[cand]
	val++
	votes[cand] = val
	return nil
}

func getMostVotes(votes map[Candidate]uint64) ([]Candidate, error) {
	var mostVotes []Candidate
	var voteCount uint64

	for cand, votes := range votes {
		if votes > voteCount {
			voteCount = votes
			mostVotes = []Candidate{cand}
			continue
		}
		if votes == voteCount {
			mostVotes = append(mostVotes, cand)
		}
	}
	return mostVotes, nil
}

type singleTransferableVote struct {
	votes map[Candidate][]Ballot
}

type candidateWithBallots struct {
	cand    Candidate
	ballots []Ballot
}

type byNumBallots []candidateWithBallots

func (a byNumBallots) Len() int           { return len(a) }
func (a byNumBallots) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byNumBallots) Less(i, j int) bool { return len(a[i].ballots) > len(a[j].ballots) }

func (stv *singleTransferableVote) Vote(b Ballot) error {
	return procBallot(stv.votes, nil, b)
}

func procBallot(
	data map[Candidate][]Ballot,
	ignores map[Candidate]struct{},
	b Ballot,
) error {
	for {
		cand, err := b.NextChoice()
		if err != nil {
			return err
		}
		if _, ok := ignores[cand]; ok {
			continue
		}
		data[cand] = append(data[cand], b)
		break
	}
	return nil
}

func (stv *singleTransferableVote) Winner() (Candidate, error) {

	localMap := map[Candidate][]Ballot{}
	for k, v := range stv.votes {
		localMap[k] = v
	}

	losers := map[Candidate]struct{}{}

	for {
		sorted := []candidateWithBallots{}
		for k, v := range localMap {
			sorted = append(sorted, candidateWithBallots{k, v})
		}
		sort.Sort(byNumBallots(sorted))

		if len(sorted) <= 2 {
			return sorted[0].cand, nil
		}

		lowest := sorted[len(sorted)-1]
		delete(localMap, lowest.cand)
		losers[lowest.cand] = struct{}{}
		for _, b := range lowest.ballots {
			err := procBallot(localMap, losers, b)
			if err != nil {
				// Ignore if you run out of candidates.
				continue
			}
		}
	}
}
