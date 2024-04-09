package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os/signal"
	"syscall"
	"word_of_wisdom/pkg"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	equihash := pkg.NewEquihashPoW(48, 3, 2)
	serv := NewServer(8080, equihash)
	go func() {
		err := serv.Listen()
		log.Panic(err)
	}()

	log.Println("server started on port 8080")
	<-ctx.Done()
	log.Println("Goodbye!")
}

type server struct {
	port     int
	equihash *pkg.EquihashPoW
}

func NewServer(port int, equihash *pkg.EquihashPoW) *server {
	return &server{
		port:     port,
		equihash: equihash,
	}
}

func (s *server) Listen() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", s.port))
	if err != nil {
		return err
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go s.handleConnection(conn)
	}
}

func (s *server) handleConnection(conn net.Conn) {
	defer conn.Close()

	challenge, err := s.equihash.Challenge()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	_, err = conn.Write(challenge)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	nonce := binary.LittleEndian.Uint64(buffer[:8])

	soln := []int{}

	for i := 8; i < n; i += 8 {
		soln = append(soln, int(binary.LittleEndian.Uint64(buffer[i:i+8])))
	}

	valid := s.equihash.Validate(int(nonce), soln, challenge)

	if !valid {
		fmt.Println("Invalid solution")
		err := conn.Close()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		return
	}

	_, err = conn.Write([]byte(wisdoms[rand.Intn(len(wisdoms))]))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

var wisdoms = []string{
	`For instance, on the planet Earth, man had always assumed that he was more intelligent than dolphins because he had achieved so much—the wheel, New York, wars and so on—whilst all the dolphins had ever done was muck about in the water having a good time. But conversely, the dolphins had always believed that they were far more intelligent than man—for precisely the same reasons.`,
	`He felt that his whole life was some kind of dream and he sometimes wondered whose it was and whether they were enjoying it.`,
	`This planet has—or rather had—a problem, which was this: most of the people living on it were unhappy for pretty much of the time. Many solutions were suggested for this problem, but most of these were largely concerned with the movement of small green pieces of paper, which was odd because on the whole it wasn't the small green pieces of paper that were unhappy.`,
	`One of the things Ford Prefect had always found hardest to understand about humans was their habit of continually stating and repeating the very very obvious.`,
	`He's spending a year dead for tax reasons.`,
	`‘Did I do anything wrong today,’ he said, ‘or has the world always been like this and I've been too wrapped up in myself to notice?’`,
	`I think you ought to know I'm feeling very depressed.`,
	`My capacity for happiness...you could fit into a matchbox without taking out the matches first.`,
	`Here, for whatever reason, is the world. And here it stays. With me on it.`,
	`Reality is frequently inaccurate.`,
	`Don't Panic.`,
	`Time is an illusion. Lunchtime doubly so.`,
	`Isn't it enough to see that a garden is beautiful without having to believe that there are fairies at the bottom of it too?`,
	`I'd far rather be happy than right any day.`,
	`If there's anything more important than my ego around, I want it caught and shot now.`,
	`‘You know,’ said Arthur, ‘it's at times like this, when I'm trapped in a Vogon airlock with a man from Betelgeuse, and about to die of asphyxiation in deep space that I really wish I'd listened to what my mother told me when I was young.’ ‘Why, what did she tell you?’ ‘I don't know, I didn't listen.’`,
	`The answer to the great question...of Life, the Universe and Everything...is...forty-two.`,
	`The argument goes something like this: ‘I refuse to prove that I exist,’ says God, ‘for proof denies faith, and without faith I am nothing.’`,
	`Anyone who is capable of getting themselves made President should on no account be allowed to do the job.`,
	`All through my life I've had this strange unaccountable feeling that something was going on in the world, something big, even sinister, and no one would tell me what it was.`,
	`Space is big. Really big. You just won't believe how vastly, hugely, mind-bogglingly big it is. I mean, you may think it's a long way down the road to the chemist's, but that's just peanuts to space.`,
	`Perhaps I'm old and tired, but I always think that the chances of finding out what really is going on are so absurdly remote that the only thing to do is to say hang the sense of it and just keep yourself occupied.`,
	`So once you do know what the question actually is, you'll know what the answer means.`,
	`Well, I mean, yes idealism, yes the dignity of pure research, yes the pursuit of truth in all its forms, but there comes a point I'm afraid where you begin to suspect that the entire multidimensional infinity of the Universe is almost certainly being run by a bunch of maniacs.`,
	`I don’t know what I’m looking for... I think it might be because if I knew I wouldn’t be able to look for them.`,
	`Looking up into the night sky is looking into infinity—distance is incomprehensible and therefore meaningless.`,
	`For a moment, nothing happened. Then, after a second or so, nothing continued to happen.`,
	`The ships hung in the sky in much the same way that bricks don't.`,
	`Ford... you're turning into a penguin. Stop it.`,
	`The last ever dolphin message was misinterpreted as a surprisingly sophisticated attempt to do a double-backwards-somersault through a hoop whilst whistling the 'Star Spangled Banner,' but in fact the message was this: ‘So long and thanks for all the fish.’`,
	`We demand rigidly defined areas of doubt and uncertainty!`,
	`What's so unpleasant about being drunk? You ask a glass of water!`,
	`In those days spirits were brave, the stakes were high, men were real men, women were real women and small furry creatures from Alpha Centauri were real small furry creatures from Alpha Centauri.`,
	`The Hitchhiker's Guide to the Galaxy also mentions alcohol. It says that the best drink in existence is the Pan Galactic Gargle Blaster. The effect of which is like having your brains smashed out with a slice of lemon wrapped round a large gold brick.`,
	`And all dared to brave unknown terrors, to do mighty deeds, to boldly split infinitives that no man had split before—and thus was the Empire forged.`,
	`Very deep... You should send that in to the Reader's Digest. They've got a page for people like you.`,
	`Why should I want to make anything up? Life’s bad enough as it is without wanting to invent any more of it.`,
	`There is an art, it says, or rather, a knack to flying. The knack lies in learning how to throw yourself at the ground and miss.`,
	`It is a mistake to think you can solve any major problems just with potatoes.`,
	`He was staring at the instruments with the air of one who is trying to convert Fahrenheit to centigrade in his head while his house is burning down.`,
	`There is a moment in every dawn when light floats, there is the possibility of magic. Creation holds its breath.`,
	`In the beginning the Universe was created. This has made a lot of people very angry and been widely regarded as a bad move.`,
}
