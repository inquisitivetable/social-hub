package seed

import (
	"SocialNetworkRestApi/api/pkg/enums"
	"time"
)

type SeedPost struct {
	Id            int
	Content       string
	PrivacyType   enums.PrivacyType
	CommentSet    []*SeedComment
	LoremComments int
	CreatedAt     time.Time
	// ImagePath   string

}

type SeedComment struct {
	UserEmail  string
	Content    string
	PostOffSet time.Duration
	// ImagePath string
}

var SeedPostsDataSetA = []*SeedPost{
	{
		Content:       "Just finished a grueling but rewarding 10K race! Who else loves the feeling of crossing the finish line?",
		PrivacyType:   enums.Public,
		CreatedAt:     time.Now().Add(time.Hour * -155),
		LoremComments: 0,
	},
	{
		Content:       "Excited to announce that my new cookbook, \"Healthy Eats for Busy Lives\", is now available for pre-order! It's packed with easy and nutritious recipes for people on the go.",
		PrivacyType:   enums.Public,
		CreatedAt:     time.Now().Add(time.Hour * -150),
		LoremComments: 0,
	},
	{
		Content:     "Who else loves practicing yoga? I just got certified as a yoga teacher and I can't wait to share my love for this practice with others.",
		PrivacyType: enums.Public,
		CreatedAt:   time.Now().Add(time.Hour * -145),
		CommentSet:  SeedCommentDataSetA1,
	},
	{
		Content:       "Happy World Vegan Day! Being vegan has been one of the best decisions I've ever made for my health and the environment.",
		PrivacyType:   enums.Public,
		CreatedAt:     time.Now().Add(time.Hour * -140),
		LoremComments: 0,
	},
	{
		Content:       "Just got back from an amazing week-long wellness retreat in Costa Rica. Yoga, healthy food, and beautiful scenery - what more could you ask for?",
		PrivacyType:   enums.Public,
		CreatedAt:     time.Now().Add(time.Hour * -135),
		LoremComments: 0,
	},
	{
		Content:       "I recently tried a new plant-based burger and it was delicious! Who says vegan food has to be boring?",
		PrivacyType:   enums.Public,
		CreatedAt:     time.Now().Add(time.Hour * -130),
		LoremComments: 0,
	},
	{
		Content:       "Who else struggles with meal planning? I'm hosting a free webinar next week on how to create a personalized meal plan that works for you. Sign up now!",
		PrivacyType:   enums.Public,
		CreatedAt:     time.Now().Add(time.Hour * -125),
		LoremComments: 0,
	},
	{
		Content:       "Just finished a challenging but rewarding HIIT workout. Who else loves a good sweat session?",
		PrivacyType:   enums.Public,
		CreatedAt:     time.Now().Add(time.Hour * -120),
		LoremComments: 0,
	},
	{
		Content:       "Happy International Women's Day! Let's celebrate the amazing women in our lives and work towards a more equal and just world.",
		PrivacyType:   enums.Public,
		CreatedAt:     time.Now().Add(time.Hour * -115),
		LoremComments: 0,
	},
	{
		Content:       "Who else loves taking their dog for a walk? My two rescue dogs, Luna and Max, are my favorite workout buddies.",
		PrivacyType:   enums.Public,
		CreatedAt:     time.Now().Add(time.Hour * -110),
		LoremComments: 0,
	},
}

var SeedPostsDataSetB = []*SeedPost{
	{
		Content:       "Just got back from an incredible trip to Japan! The food, the people, and the culture were all amazing. Can't wait to go back someday!",
		PrivacyType:   enums.Public,
		CreatedAt:     time.Now().Add(time.Hour * -105),
		LoremComments: 0,
	},
	{
		Content:       "Who else loves a good outdoor adventure? Just went on a challenging hike up Mount Kilimanjaro and it was definitely worth it. What's your favorite hike?",
		PrivacyType:   enums.Public,
		CreatedAt:     time.Now().Add(time.Hour * -100),
		LoremComments: 0,
	},
	{
		Content:       "Excited to share that my latest travel article was published in National Geographic! It's all about the hidden gems of Barcelona. Check it out if you're planning a trip there soon!",
		PrivacyType:   enums.Public,
		CreatedAt:     time.Now().Add(time.Hour * -95),
		LoremComments: 0,
	},
	{
		Content:       "Recently started a photography project where I take a photo of something beautiful every day. It's been a great way to appreciate the little things in life.",
		PrivacyType:   enums.Public,
		CreatedAt:     time.Now().Add(time.Hour * -90),
		LoremComments: 0,
	},
	{
		Content:       "Happy National Book Lovers Day! I just finished \"The Overstory\" by Richard Powers and it's definitely one of my new favorites. What are you currently reading?",
		PrivacyType:   enums.Public,
		CreatedAt:     time.Now().Add(time.Hour * -85),
		LoremComments: 0,
	},
	{
		Content:       "I'm officially a certified scuba diver! I've always been fascinated by the ocean and it was incredible to see all the marine life up close.",
		PrivacyType:   enums.Public,
		CreatedAt:     time.Now().Add(time.Hour * -80),
		LoremComments: 0,
	},
	{
		Content:       "Anyone else a fan of street food? I recently tried the best tacos al pastor I've ever had in Mexico City. Already planning my next trip back!",
		PrivacyType:   enums.Public,
		CreatedAt:     time.Now().Add(time.Hour * -75),
		LoremComments: 0,
	},
	{
		Content:     "Just got back from an amazing trip to Bali. The beaches, the temples, and the food were all incredible. Can't wait to go back someday!",
		PrivacyType: enums.Public,
		CreatedAt:   time.Now().Add(time.Hour * -70),

		LoremComments: 0,
	},
	{
		Content:       "Just finished a 10-day silent meditation retreat and it was one of the most challenging and rewarding experiences of my life. Highly recommend it to anyone interested in mindfulness and inner peace.",
		PrivacyType:   enums.Public,
		CreatedAt:     time.Now().Add(time.Hour * -65),
		LoremComments: 0,
	},
	{
		Content:       "Happy Earth Day! Let's all do our part to protect our planet and make it a better place for future generations.",
		PrivacyType:   enums.Public,
		CreatedAt:     time.Now().Add(time.Hour * -60),
		LoremComments: 0,
	},
}

var SeedPostsDataSetC = []*SeedPost{
	{
		Content:       "Just launched my latest app, \"TaskMaster\", on the App Store! It's a productivity app that helps you stay on top of your to-do list. Check it out!",
		PrivacyType:   enums.Public,
		CreatedAt:     time.Now().Add(time.Hour * -55),
		LoremComments: 0,
	},
	{
		Content:       "Who else loves a good hackathon? Just won first place at the HackNY hackathon with my team. Can't wait for the next one.",
		PrivacyType:   enums.Public,
		CreatedAt:     time.Now().Add(time.Hour * -50),
		LoremComments: 2,
	},
	{
		Content:       "Excited to announce that I've been accepted into the Google Developer Expert program for mobile development! It's an honor to be part of this community of experts.",
		PrivacyType:   enums.Public,
		CreatedAt:     time.Now().Add(time.Hour * -45),
		LoremComments: 3,
	},
	{
		Content:       "Just wrapped up a successful project with a Fortune 500 company. It was a challenging but rewarding experience, and I'm proud of what our team accomplished.",
		PrivacyType:   enums.Public,
		CreatedAt:     time.Now().Add(time.Hour * -40),
		LoremComments: 0,
	},
	{
		Content:       "Excited to share that I'll be speaking at the upcoming TechCrunch Disrupt conference about the future of mobile development. Can't wait to share my insights with the tech community!",
		PrivacyType:   enums.Public,
		CreatedAt:     time.Now().Add(time.Hour * -35),
		LoremComments: 22,
	},
	{
		Content:       "Feeling a bit burnt out lately. The tech industry can be so demanding sometimes, and I feel like I'm always on call. Trying to take some time for self-care and relaxation, but it's tough when there's always another deadline looming.",
		PrivacyType:   enums.Private,
		CreatedAt:     time.Now().Add(time.Hour * -30),
		LoremComments: 0,
	},
	{
		Content:       "Had a heart-to-heart with my mentor today about imposter syndrome. It's something that's been weighing on me lately, but it was helpful to hear that even seasoned developers experience it from time to time.",
		PrivacyType:   enums.Private,
		CreatedAt:     time.Now().Add(time.Hour * -25),
		LoremComments: 0,
	},
	{
		Content:       "Dealing with a bit of imposter syndrome lately. I keep worrying that I'm not skilled enough or experienced enough to tackle the projects I'm working on. Trying to remind myself that I wouldn't have gotten this far if I didn't have the skills and knowledge to back it up.",
		PrivacyType:   enums.Private,
		CreatedAt:     time.Now().Add(time.Hour * -20),
		LoremComments: 0,
	},
	{
		Content:       "Just got back from a weekend getaway with my partner. It was so nice to unplug from work and spend some quality time together.",
		PrivacyType:   enums.Private,
		CreatedAt:     time.Now().Add(time.Hour * -15),
		LoremComments: 0,
	},
	{
		Content:       "Feeling grateful for my team today. We've been working on a really challenging project, but everyone has been pulling their weight and pushing us towards success. It's great to work with such talented and dedicated individuals.",
		PrivacyType:   enums.Private,
		CreatedAt:     time.Now().Add(time.Hour * -10),
		LoremComments: 0,
	},
}

// 6 rows of ChatGPT generated comments for user A
var SeedCommentDataSetA1 = []*SeedComment{
	{
		UserEmail:  "b@b.com",
		Content:    " Congrats on getting certified, that's awesome! Can't wait to attend one of your classes and learn from the best.",
		PostOffSet: time.Minute * 360,
	},
	{
		UserEmail:  "a@a.com",
		Content:    "@Benjamin, Thank you!",
		PostOffSet: time.Minute * 330,
	},
	{
		UserEmail:  "c@c.com",
		Content:    "I'm a big fan of yoga too! It's such a great way to unwind and de-stress after a long day at work. Congrats on becoming a teacher!",
		PostOffSet: time.Minute * 270,
	},
	{
		UserEmail:  "a@a.com",
		Content:    "@Carlos, yoga has been a lifesaver for me in terms of managing my stress and anxiety, and I'm sure you'll find it helpful too.",
		PostOffSet: time.Minute * 240,
	},
	{
		UserEmail:  "d@d.com",
		Content:    "I've been wanting to try yoga for ages, but never found the right teacher. Looking forward to attending one of your classes and finally giving it a go!",
		PostOffSet: time.Minute * 180,
	},
	{
		UserEmail:  "a@a.com",
		Content:    "@Deanna, I was in the same boat as you before I found the right teacher - trust me, it's worth the wait! Looking forward to seeing you all at the studio soon.",
		PostOffSet: time.Minute * 150,
	},
	{
		UserEmail:  "e@e.com",
		Content:    "Yoga is a game-changer! I've been practicing for years and it's helped me maintain both my physical and mental health. Excited to see you share your knowledge and passion with others.",
		PostOffSet: time.Minute * 90,
	},
}
