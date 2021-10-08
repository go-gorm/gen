package check

type MemberOpt interface {
	self() func(*Member) *Member
}

type ModifyMemberOpt func(*Member) *Member

func (o ModifyMemberOpt) self() func(*Member) *Member { return o }

type FilterMemberOpt ModifyMemberOpt

func (o FilterMemberOpt) self() func(*Member) *Member { return o }

type CreateMemberOpt ModifyMemberOpt

func (o CreateMemberOpt) self() func(*Member) *Member { return o }

func sortOpt(opts []MemberOpt) (modifyOpts []MemberOpt, filterOpts []MemberOpt, createOpts []MemberOpt) {
	for _, opt := range opts {
		switch opt.(type) {
		case ModifyMemberOpt:
			modifyOpts = append(modifyOpts, opt)
		case FilterMemberOpt:
			filterOpts = append(filterOpts, opt)
		case CreateMemberOpt:
			createOpts = append(createOpts, opt)
		}
	}
	return
}
