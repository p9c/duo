#### Parallelcoin Technology Brief 12 August 2018

# Reputation

> This document outlines the problem domain for which the projects of the Parallelcoin Team are working to solve

Reputation is the esteem given to an individual or organisation (indirectly, but opaquely the responsible individuals) towards a particular quality or more especially, propriety.

In the real world, there is overt and latent reputation records, some declarations are recorded and certified, some are not. People may refuse to even speak let alone sign off on an account of a complaint (or praise) on the basis of some kind of belief.

There will always of course be the problem of unexpressed (dis)approvals, but lowering the labor cost of issuing them would afford a great deal of information to benefit of all individuals of the data of the evaluation of as many as possible of the counterparties a prospective trading partner might have been involved with.

So there is also the problem of attribution. How many more, *at least equipotentially positive*, (dis)approvals would be recorded if the authors could be certain that they would not be attributable yet at the same time be verifiable.

The only way to make this both fair and quantifiable, is to continuously compute the flux of approval and disapproval between users by treating it as a type of currency. The interface is arbitrary, and probably should be user-configurable.

An account has some amount of reputation currency, and this currency balance influences the computation of rewards for other actions taken by an account holder and certified by necessary third parties via secure protocols.

Through the use of "Confidential Transactions" it is possible to spend reputation, create anonymised (which are also reduced by some percentage in effect, also compensating for the cryptography required) complaints and praise, even to register shunning someone without it being known it is you, without being able to cheat the amount of effect you can have.

## Types of reputation spends

### Vote

Votes are like a payment. The amount of reputation spent, minus consensus based commissions for various purposes, is part of a computation that issues an amount of currency of the main, transactional currency, based on these social network (dis)approval actions.

### Hold

Holds are enacted as part of a declaration to the network that you do not wish to see the posts of some user. These declarations have the effect, in proportion with one's balance of reputation, of placing a proportional hold on the target. Thus both hater and hated lose the ability to spend some of their reputation.

The measure is temporary, effective until canceled or potentially expired or decayed at some prescribed rate.

## Importance of possible anonymity

While it might be said that anonymous complaints can cause a problem, actually, the real problem is that certified, but not signed statements, might encourage lynch mob behaviour. But there is another side to this. Sometimes fear motivates an ongoing reticence towards being credited for disapproving, and a number of anonymous complaints might sometimes break the silence.

Signatures are worth something. They are worth as much as what they are put upon if the document is valid (eg, account has enough to pay out a cheque). So the cost of anonymous (dis)approvals is a significant reduction in their effect.

People are more prone to flamboyance behind a screen, so in itself, a transparent record of reputation has some kind of value, as it would shift the balance further towards the centre against those with the power to punish complainers.

Secondly, it is possible, at a cost of discounted reputation, to register approval/disapproval without declaring this publicly.

## The forum as the centre of reputation discovery

Commerce is a very broadly applicable word. In essence, in most dictionaries, it means any kind of social interaction. In commerce there is specialisations. Thus a user would have several reputation scores for each of the forums they interact with. A forum is an environment in which commerce occurs.

Forums include:

1. Bulletin Board Systems
2. Auction and Marketplace Systems
3. Peer review data repositories (encyclopaedia, research publication, open source software repositories)

...and many more things also - though the specifics might vary the general character of being commerce will be present.

It covers everything. Someone gives you stellar table service, you can spend a big whack of your rep balance on them. The transaction will finalise a reward a week later based on the amount of reputation spent.

## The problem of Sybil Attacks

Social networks, being a similar category and architecture to democracy, are vulnerable to deliberate distortion of the population, vote stuffing, and the social network equivalent, a second, intentionally dissociated identity through which attacks are launched.

The only reliable method so far found for discovering sybils in a social network is the efforts of usually many individual users to investigate, be it simply analysis or more intrusive attacks aimed at unmasking the suspected clone account.

### Crowd-sourced security against Sybils

Thus, the most effective way to suppress sybil, sock-puppet activity, is to empower *all* users to have the ability to apply a reduction on the ranking of a suspected bad actor and to prevent this from being abused, the effect of negative votes and holds is equally suffered by both accounts.

The compounding of numbers of complaining accounts, firstly being limited by the account's balance of reputation, and secondly by a network consensus rule that in (network consensus) chronological order, each subsequent user's down vote or hold is diminished at a rate that keeps the size of an interniecine conflict to about 60 individuals, as more users contributing together beyond this level have insignificant effect on the outcome.

## Upvotes, Downvotes, Shunning and Following

Each of these actions are transactions in the distributed database. Each one has an effect based on associated actions related to such things as posts (any media).

### **Upvotes** simply assign, like a payment, the reputation credit to the target.

The total reward of votes is computed and scheduled to be processed in the 24 hour period of the 7th day after the original posting.

It is network consensus that one may not vote for oneself. It is like wanting the public to acknowledge you gave yourself money you already owned. The prevention of use of Sybils to attempt to manipulate the voting outcome is as already discussed, not protocol level but for this reason has a negative way to use the reputation credits.

So rather than delve into the sophistic arguments that 'sybils cannot be eliminated' and thus must be 'tolerated' - no, through the mechanism of distribution of reputation credits, the acquisition of serial bad actors of unearned reputation can be countered by the burning of reputation by those who have fairly earned it.

And it is hardly new - in the Code of Hammurabi, false accusations proven to be false are subject to the identical punishment as the alleged crime. So in this system, according to your feelings about it, you can dump any or all of your earned reputation on burying someone. But actually the math of the system is such that even with a big crowd of reputable users the effect is capped so that 'unpersoning' means losing potentially as much as 10% from one heavily downvoted post, or at 50% limit on power from shunning. The proportioning is in order to allow people to learn how to use it their own way. Ultimately how this proposed 'currency' will function will entirely depend on its suitability and marketability.

In these current times, it is unquestionable that a clear record of cumulated evaluations by people operating through digital intermediaries are more likely to be truthful about (especially if it costs them the power to assign value to another's post and pool of reputation to spend), as well as preventing the total silencing of a user's posts. We don't care if you want to teach your pug to do a nazi/roman salute. We just probably don't want to see it unless we find such a thing comical. We can use also the publicly signed off feedback evaluations of others towards others to filter exactly how we decide to deal with another user. Maybe no matter that 60% think some person is a fountain of truth and beauty, you might think they are crap and the network client enables the user to filter it this way.

### Downvotes also spend reputation credits but spend, proportionally, the reputation credits of the account making a downvoted post

The downvote also results in the burning of the downvoter's reputation credits, as well as a progressively diminished burn on the target account's reputation credit balance, which cancels out some of the reward that would be added for the upvotes received.

*The reduction of downvote effect as numbers increase (in accounts) is for two reasons:*

1. Sybils can pump numbers. They can only be eliminated by lowering their reward from exploiting users to the point that honest behaviour is more profitable.

2. In almost all systems of law, and custom, it is considered the magic number 3 for proof of recidivism. Thus, if at most a user can lose half of their reputation in one week, they get 3 weeks to turn it around and in absolute numbers the reduction of influence caused diminishes against the previous time. 50% per strike means 75% from two and 87.5% from strike three. In terms of theoretical maximum burn of reputation. Plus a similarly count-reduced effect (and order effect) of holds can make a person nearly a non-person in terms of whether anyone cares if they are downvoted or shunned by them.

### Requiring Infrastructure Providers to be Accountable

The most central individuals in a distributed, cryptographically secured database network system are the ones who run the machines that relay and process and certify the data for other users. Indeed, in the internet in general. They necessarily must have a higher standard of ethics or they can cause their fellow users to diminish their operators reward to zero or even a cost.

For this reason, the default trust score, in the general forum, influences the default selection frequency of a given account to provide monetisable services. In a poetic way, faithful relaying of speech implies a commitment to being truthful, so if a prominent infrastructure provider goes mustang, maybe they are positively implicated in an attempt to attack the network for profit, or for any reason at all, they can have their pay rate slashed and thereby brought into line or to settle the case in their favour should they be in the right.

Infrastructure operators are like representatives. They literally re-present the data you send to them. Thus they must be transparent in their methods. This is vital to trust in the network and most especially in the utility and efficacy of the reputation credit system.

### Dealing with complaints is very important

It should advantage nobody for any other factor than their veracity and innocence when a conflict arises. It should not be potentially profitable to complain in exchange for money. It should not be possible for one individual to exert the effect of many upon the reputation of another person without a serious cost to the recreational whiner.

And at the same time, though it should be discounted to some degree, anonymous complaints should have a real, tangible influence on the target, for the reason that these weakened complaints act like a semaphore. Not a strong effect but enough to get attention, enough to be worth paying the premium to file it.

## Reputation as a currency, but with additional rules of exchange

A normal money does not have the property of affecting the accounts of other people, except in as far as executing a transfer of value that first came into the account legitimately.

A reputation currency is like a normal money, except that to some extent in proportion with your holdings, you can make other people's tokens also disappear. Its power is in quantifying socially consensual agreement as to the valuation of a piece of work.

"Negative" votes have the effect of reducing currency supply. The currency is issued to accounts from a baseline minimum based on the 'starter average' level and zero issuance occurs towards an account while its reputation balance is negative and the same for a downvoted post. Thus, those who are motivated to engage in policing behaviour, do so at the cost of their ability to reward other's work and by shrinking the supply, differentially give more reward to those who don't engage in 'negative' transactions.

## Conclusion

The purpose of all systems of law is to minimise conflict and the damage that can result while it remains unresolved. The perfect solution is having everything documented, but many things cannot be documented. It should be at the prerogative and expense of an individual to disavow identification, as this is still a greater social good than no feedback at all.

The behaviour which such a system engenders is largely positive, towards pleasing more people with your ramblings and masterpieces, and reducing the amount of those who try to advocate violence without the deserved reduction of reputation such acts declarations should attract. It functions as an effective distribution system for compensating artists, who are by their nature generally basically beggars, sometimes with airs.

Rather than attempt the impossible and stop Sybil attacks, empower the victims to mobilise and counter such mischief.

It is indeed a social scoring system and it is admittedly narrow in its scope in that it evaluates public declarations of various kinds.

However, by scrupulously quantifying the value of these evaluations, it becomes possible to more gently persuade people towards good behaviour, and it eliminates the problem of a corrupted establishment intermediary. Enough users with sufficient reputation, it doesn't matter where it comes from, can halt the undesired behaviour of individuals on the network. It becomes a kind of court, and thus, as we are designing the court stenographer to machines, saving the fun work for those who are interested and capable of exerting pressure towards a positive outcome.