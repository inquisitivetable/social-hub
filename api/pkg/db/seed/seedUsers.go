package seed

type SeedUser struct {
	Id              int64
	FirstName       string
	LastName        string
	Email           string
	Nickname        string
	About           string
	ImagePath       string
	IsPublic        bool
	PostSet         []*SeedPost
	FollowingEmails []string
}

var SeedUserData = []*SeedUser{
	{
		FirstName:       "Ann",
		LastName:        "Addams",
		Email:           "a@a.com",
		Nickname:        "AnnieA",
		About:           "Hey there, I'm Ann! I'm a fitness enthusiast and nutrition coach. I'm passionate about helping others achieve their health and wellness goals through personalized meal plans and exercise routines. In my free time, I enjoy trying out new healthy recipes, practicing yoga, and spending time with my two rescue dogs.",
		IsPublic:        true,
		PostSet:         SeedPostsDataSetA,
		FollowingEmails: []string{"c@c.com"},
		ImagePath:       "UserA.jpg",
	},
	{
		FirstName:       "Benjamin",
		LastName:        "Button",
		Email:           "b@b.com",
		Nickname:        "Buttons",
		About:           "Hi, I'm Benjamin! I'm a freelance writer and avid traveler. I love exploring new cultures and cuisines, and I'm always on the lookout for my next adventure. When I'm not writing or traveling, you can find me hiking, reading a good book, or practicing my photography skills.",
		IsPublic:        true,
		PostSet:         SeedPostsDataSetB,
		FollowingEmails: []string{"a@a.com", "f@f.com"},
		ImagePath:       "UserB.jpg",
	},
	{
		FirstName: "Carlos",
		LastName:  "Cortez",
		Email:     "c@c.com",
		Nickname:  "Carlito",
		About:     "Hi, my name is Carlos! I'm a software developer with a love for all things tech. I specialize in building mobile and web applications, and I'm always looking for new and innovative ways to solve complex problems through code. When I'm not coding, you can find me playing video games or tinkering with my latest DIY project.",
		IsPublic:  false,
		PostSet:   SeedPostsDataSetC,
		ImagePath: "UserC.jpg",
	},
	{
		FirstName: "Deanna",
		LastName:  "Davis",
		Email:     "d@d.com",
		Nickname:  "DeeDee",
		About:     "Hi there, I'm Deanna Davis. I'm a freelance writer and digital marketer with a passion for creating compelling content that connects with audiences. When I'm not working, you can usually find me hiking with my dog or experimenting with new vegan recipes in the kitchen.",
		IsPublic:  true,
		ImagePath: "UserD.jpg",
	},
	{
		FirstName: "Ethan",
		LastName:  "Evans",
		Email:     "e@e.com",
		Nickname:  "EvansCode",
		About:     "Hey, I'm Ethan Evans. I'm a software engineer with over 10 years of experience in the industry. I'm passionate about using technology to solve real-world problems and improve people's lives. When I'm not coding, I enjoy playing basketball and exploring new restaurants in the city.",
		IsPublic:  true,
	},
	{
		FirstName:       "Felix",
		LastName:        "Fitzgerald",
		Email:           "f@f.com",
		Nickname:        "FelixTheExplorer",
		About:           "Hi, I'm Felix! I'm here to bring the fun to your screen with my adventures, jokes, and love for all things entertaining. Let's have a blast together!",
		FollowingEmails: []string{"a@a.com", "b@b.com", "c@c.com", "d@d.com", "e@e.com"},
		IsPublic:        true,
	},
	{
		FirstName:       "Grace",
		LastName:        "Garner",
		Email:           "g@g.com",
		Nickname:        "GracefulG",
		About:           "Hey there! I'm Grace, and I'm all about spreading good vibes. Join me on my groovy journey filled with positivity, music, and inspiring moments.",
		FollowingEmails: []string{"c@c.com", "f@f.com"},
		IsPublic:        false,
		ImagePath:       "UserG.jpg",
	},
	{
		FirstName:       "Harper",
		LastName:        "Harrison",
		Email:           "h@h.com",
		Nickname:        "Harrison Hord",
		About:           "Welcome to my world of happiness! I'm Harper, a joyful soul who loves to explore, laugh, and share smiles. Let's create happy memories together!",
		FollowingEmails: []string{"c@c.com", "f@f.com"},
		IsPublic:        true,
	},
	{
		FirstName:       "Isaac",
		LastName:        "Ingram",
		Email:           "i@i.com",
		Nickname:        "ImpressiveIsaac",
		About:           "Greetings, folks! I'm Isaac, and I'm here to impress you with my musical talents, random facts, and a sprinkle of charm. Let's make some magic!",
		FollowingEmails: []string{"f@f.com"},
		IsPublic:        false,
	},
	{
		FirstName: "Jasmine",
		LastName:  "Jenkins",
		Email:     "j@j.com",
		Nickname:  "10Jasmine01",
		About:     "Embrace the joy! I'm Jasmine, a foodie, traveler, and happiness seeker. Join me on a journey filled with flavors, laughter, and beautiful moments.",
		IsPublic:  true,
	},
	{
		FirstName: "Lorem",
		LastName:  "Ipsum",
		Email:     "l@l.com",
		Nickname:  "LoremGenerator",
		About:     "",
	},
}
