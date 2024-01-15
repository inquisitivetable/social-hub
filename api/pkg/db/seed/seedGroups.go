package seed

import (
	"time"
)

type SeedGroup struct {
	CreatorEmail       string
	Title              string
	Description        string
	Users              []string
	ImagePath          string
	SeedEventsData     []*SeedEvent
	SeedGroupPostsData []*SeedGroupPost
}

type SeedGroupPost struct {
	Id           int
	CreatorEmail string
	Content      string
	// CommentSet    []*SeedComment
	LoremComments int
	CreatedAt     time.Time
}

var SeedGroupsData = []*SeedGroup{
	{
		CreatorEmail:       "b@b.com",
		Title:              "Adventurers United",
		Description:        "Adventurers United is a group dedicated to exploring uncharted territories, unraveling mysteries, and pushing the boundaries of human exploration. Comprising fearless individuals from different backgrounds, this group seeks to discover hidden wonders, encounter diverse cultures, and document their findings for the world to marvel at. From scaling towering peaks to delving into ancient ruins, Adventurers United is fueled by the passion for discovery and the thrill of venturing into the unknown.",
		Users:              []string{"a@a.com", "b@b.com", "c@c.com", "d@d.com", "e@e.com"},
		SeedEventsData:     SeedEventsDataA,
		SeedGroupPostsData: SeedGroupPostsDataSetA,
		ImagePath:          "groupA.png",
	},
	{
		CreatorEmail: "a@a.com",
		Title:        "Blissful Harmony",
		Description:  "Blissful Harmony is a group dedicated to promoting inner peace and cultivating a harmonious existence through mindfulness, meditation, and holistic practices.",
		Users:        []string{"a@a.com", "d@d.com", "e@e.com"},
		ImagePath:    "groupB.png",
	},
	{
		CreatorEmail:   "b@b.com",
		Title:          "Creative Catalysts",
		Description:    "Creative Catalysts is a gathering of visionary artists, innovators, and thinkers who believe in the transformative power of creativity. Embracing diverse forms of artistic expression, this group seeks to challenge conventions, provoke thought, and inspire change through their work. From visual arts and music to literature and performance, Creative Catalysts use their talents to ignite conversations, bridge gaps, and breathe life into new ideas. Their collective energy fuels a vibrant and dynamic creative community.",
		Users:          []string{"a@a.com", "b@b.com", "c@c.com", "d@d.com"},
		SeedEventsData: SeedEventsDataC,
		ImagePath:      "groupC.png",
	},
	{
		CreatorEmail:       "c@c.com",
		Title:              "Dreamcatchers",
		Description:        "Dreamcatchers is a group dedicated to empowering individuals to pursue their dreams and aspirations. Recognizing the importance of nurturing ambitions, this group provides a supportive platform for members to share their goals, seek guidance, and celebrate achievements. Through mentorship programs, motivational workshops, and networking opportunities, Dreamcatchers aim to inspire personal growth, foster resilience, and create a community of dreamers who believe that anything is possible with dedication and perseverance.",
		Users:              []string{"c@c.com", "d@d.com", "e@e.com"},
		SeedGroupPostsData: SeedGroupPostsDataSetD,
		ImagePath:          "groupD.png",
	},
	{
		CreatorEmail: "b@b.com",
		Title:        "Ecological Guardians",
		Description:  "Ecological Guardians is a passionate group of environmental advocates committed to preserving and protecting the planet. Driven by a deep concern for the Earth's well-being, they actively engage in conservation efforts, sustainable practices, and education campaigns to raise awareness about environmental issues. From organizing clean-up drives to promoting renewable energy solutions, Ecological Guardians work tirelessly to safeguard ecosystems, combat climate change, and inspire others to adopt eco-friendly lifestyles for a greener future.",
		Users:        []string{"b@b.com"},
		ImagePath:    "groupE.png",
	},
	{
		CreatorEmail: "h@h.com",
		Title:        "Friends of Nature",
		Description:  "Friends of Nature is a community of nature enthusiasts who are passionate about exploring and preserving the natural world. Through their shared love for the environment, they engage in various outdoor activities, such as hiking, wildlife observation, and nature photography. Members of Friends of Nature also actively participate in conservation projects, organize awareness campaigns, and advocate for sustainable practices to protect ecosystems and biodiversity.",
		Users:        []string{"a@a.com", "b@b.com", "c@c.com", "d@d.com", "e@e.com", "h@h.com"},
		ImagePath:    "groupF.png",
	},
	{
		CreatorEmail: "h@h.com",
		Title:        "Global Changemakers",
		Description:  "Global Changemakers is a diverse community of individuals dedicated to making a positive impact on the world. The group focuses on addressing various global challenges, including poverty, inequality, and environmental degradation. Members of Global Changemakers collaborate on social initiatives, volunteer in local communities, and support sustainable development projects. Through their collective efforts, they strive to create lasting change and build a more equitable and sustainable future for all.",
		Users:        []string{"h@h.com"},
		ImagePath:    "",
	},
	{
		CreatorEmail: "a@a.com",
		Title:        "Health and Wellness Enthusiasts",
		Description:  "Health and Wellness Enthusiasts is a vibrant group of individuals passionate about promoting holistic well-being. The group focuses on sharing knowledge, tips, and resources related to physical fitness, mental health, nutrition, and self-care practices. Members of Health and Wellness Enthusiasts engage in discussions, organize fitness challenges, and support each other on their wellness journeys. By fostering a supportive and inclusive community, the group aims to inspire others to prioritize their health and embrace a balanced lifestyle.",
		Users:        []string{"a@a.com", "d@d.com", "e@e.com"},
		ImagePath:    "groupH.png",
	},
	{
		CreatorEmail: "c@c.com",
		Title:        "Innovation Hub",
		Description:  "Innovation Hub is a collaborative space for creative thinkers, entrepreneurs, and innovators to share ideas and explore new possibilities. The group encourages discussions on emerging technologies, startup ventures, and breakthrough innovations across various industries. Members of the Innovation Hub network, exchange expertise, provide feedback, and foster a supportive environment for turning ideas into reality. With a shared passion for innovation, the group aims to inspire and drive positive change in the world through entrepreneurial endeavors.",
		Users:        []string{"c@c.com", "e@e.com"},
		ImagePath:    "",
	},
	{
		CreatorEmail: "j@j.com",
		Title:        "Joyful Journeyers",
		Description:  "Joyful Journeyers is a community of travel enthusiasts who believe in the transformative power of exploration and cultural immersion. The group shares travel experiences, tips, and recommendations to inspire others to embark on enriching journeys. Members of Joyful Journeyers engage in conversations about different destinations, local customs, and sustainable travel practices. They also organize group trips and collaborate on travel-related projects aimed at promoting responsible tourism and cross-cultural understanding. Joyful Journeyers encourages members to embrace the joy of travel and create unforgettable experiences.",
		Users:        []string{"b@b.com", "j@j.com", "f@f.com", "h@h.com"},
		ImagePath:    "",
	},
}

var SeedGroupPostsDataSetA = []*SeedGroupPost{
	{
		Content:       "Are you passionate about exploring uncharted territories, uncovering ancient mysteries, and pushing the boundaries of human exploration? Look no further! Adventurers United welcomes fearless individuals from all walks of life to join us on our quest for hidden wonders and diverse cultures. Together, let's document our thrilling findings and inspire the world with the marvels of the unknown! 🧭🌍 #AdventurersUnited #ExploreTheUnknown #DiscoverWonders",
		CreatedAt:     time.Now().Add(time.Hour * -195),
		CreatorEmail:  "b@b.com",
		LoremComments: 5,
	},
	{
		Content:       "As a member of Adventurers United, I am thrilled to be part of a fearless community that explores uncharted territories and unravels mysteries hidden in the depths of history. From scaling towering peaks to delving into ancient ruins, every expedition brings us closer to the wonders of our world. Join us on this exhilarating journey as we push the boundaries of human exploration and celebrate the thrill of venturing into the unknown! #AdventurersUnited #ExplorationThrills #DiscoverTheUnknown",
		CreatedAt:     time.Now().Add(time.Hour * -190),
		CreatorEmail:  "a@a.com",
		LoremComments: 10,
	},
	{
		Content:       "Adventurers United is not just a group of explorers; we are a diverse tapestry of backgrounds and cultures, united by our shared passion for discovery. Together, we encounter diverse cultures and embrace the beauty of our global family. The world is our playground, and we document our findings for the world to marvel at. Come, join our tribe, and embark on a journey that celebrates the wonders of our planet and the spirit of unity in exploration!  #AdventurersUnited #EmbraceDiversity #GlobalExplorers",
		CreatedAt:     time.Now().Add(time.Hour * -185),
		CreatorEmail:  "c@c.com",
		LoremComments: 11,
	},
	{
		Content:       "Adventurers United is more than just a group; it's a way of life. We seek out uncharted territories and scale towering peaks, not just for the thrill but to test our limits and find inner strength. With every expedition, we challenge ourselves and come out stronger. We carry with us the spirit of exploration, empowering ourselves to face the unknown with determination and courage. Adventure awaits, and we invite you to join our tribe of explorers, breaking barriers and embracing the thrill of the journey! #AdventurersUnited #ScalingNewHeights #ThrillOfExploration		",
		CreatedAt:     time.Now().Add(time.Hour * -180),
		CreatorEmail:  "d@d.com",
		LoremComments: 1,
	},
	{
		Content:       "As a member of Adventurers United, I am honored to be part of a group that not only seeks out hidden wonders but also documents our expeditions to share with the world. Through our photographs, stories, and experiences, we aim to inspire others to embrace the unknown and appreciate the diverse beauty of our planet. Our journey is not just about personal growth; it's about leaving a mark on the world, a legacy of exploration that future generations can cherish. Together, we make history as we venture into the uncharted, one step at a time. #AdventurersUnited #DocumentingExpeditions #LeaveYourLegacy		",
		CreatedAt:     time.Now().Add(time.Hour * -170),
		CreatorEmail:  "a@a.com",
		LoremComments: 0,
	},
}

var SeedGroupPostsDataSetD = []*SeedGroupPost{
	{
		Content:       "I can't express how amazing it is to be surrounded by fellow dreamers who believe in the power of aspirations. The support and encouragement I've received here have given me the confidence to pursue my dreams fearlessly. From sharing my goals to attending motivational workshops, every step has been instrumental in my personal growth journey. Thank you, Dreamcatchers, for empowering me to believe that anything is possible with dedication and perseverance! #Dreamcatchers #EmpoweringDreams #BelieveInYourself",
		CreatedAt:     time.Now().Add(time.Hour * -175),
		CreatorEmail:  "c@c.com",
		LoremComments: 5,
	},
	{
		Content:       "Today, I reached a significant goal I set for myself, and it fills my heart with joy to share this victory with all of you! Dreamcatchers has been an unwavering source of motivation and support throughout this journey. From the insightful workshops to the kind words of encouragement, this community has been a driving force behind my progress. To all fellow dreamers, keep pushing forward – together, we prove that anything is achievable with dedication and perseverance! #Dreamcatchers #AchievementUnlocked #DreamChaser",
		CreatedAt:     time.Now().Add(time.Hour * -170),
		CreatorEmail:  "c@c.com",
		LoremComments: 10,
	},
	{
		Content:       "Being a part of Dreamcatchers has introduced me to some incredible mentors who have selflessly shared their wisdom and experiences. Their guidance has been invaluable in shaping my path towards achieving my dreams. I've learned that resilience and determination are key to overcoming obstacles, and I'm excited to keep striving for my goals with the support of this amazing community. Dreamcatchers, you truly inspire personal growth and make dreams possible! #Dreamcatchers #MentorshipProgram #Empowerment",
		CreatedAt:     time.Now().Add(time.Hour * -165),
		CreatorEmail:  "d@d.com",
		LoremComments: 11,
	},
	{
		Content:       "Life's journey may not always be smooth sailing, but with the Dreamcatchers community by my side, I've learned to embrace challenges as opportunities for growth. Together, we share our dreams, support each other during tough times, and celebrate the victories that come our way. It's a constant reminder that no dream is too big and that dedication and perseverance are the keys to success. Grateful for this inspiring group that fosters resilience and empowers us all to reach for the stars! #Dreamcatchers #Resilience #DreamBig",
		CreatedAt:     time.Now().Add(time.Hour * -160),
		CreatorEmail:  "d@d.com",
		LoremComments: 1,
	},
}
