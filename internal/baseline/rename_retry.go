package baseline

import "time"

func retryRename(rename func() error, isEphemeral func(error) bool, sleep func(time.Duration), limit time.Duration) error {
	delay := time.Millisecond
	var elapsed time.Duration
	for {
		err := rename()
		if err == nil || !isEphemeral(err) {
			return err
		}
		if elapsed >= limit {
			return err
		}

		wait := delay
		if remaining := limit - elapsed; wait > remaining {
			wait = remaining
		}
		sleep(wait)
		elapsed += wait

		delay *= 2
		if delay > 100*time.Millisecond {
			delay = 100 * time.Millisecond
		}
	}
}
