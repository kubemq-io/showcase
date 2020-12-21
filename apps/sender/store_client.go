package main

//
//import (
//	"context"
//	"fmt"
//	"github.com/kubemq-io/kubemq-go"
//	"github.com/nats-io/nuid"
//	"log"
//
//	"time"
//)
//
//type StoreClient struct {
//	Id            int
//	cfg           *Config
//	stats         *ClientStats
//	localClientID string
//	localChannel  string
//	storage       *Storage
//}
//
//func NewStoreClient(ctx context.Context, id int, cfg *Config) *StoreClient {
//	c := &StoreClient{
//		Id:            id,
//		cfg:           cfg,
//		stats:         NewClientStats(),
//		localClientID: fmt.Sprintf("%s-%d", cfg.ClientId, id),
//		localChannel:  fmt.Sprintf("%s.%d", cfg.Channel, id),
//		storage:       nil,
//	}
//	c.storage = InitStorage(c.localChannel)
//	go c.RunSenders(ctx)
//	return c
//}
//func (c *StoreClient) Log(msg string) {
//	if c.cfg.Verbose {
//		log.Println(msg)
//	}
//}
//func (c *StoreClient) Logf(format string, args ...interface{}) {
//	if c.cfg.Verbose {
//		log.Println(fmt.Sprintf(format, args...))
//	}
//}
//
//func (c *StoreClient) RunSenders(ctx context.Context) {
//	randomMessageBody := randomBytes(c.cfg.PayloadSize)
//	for i := 1; i <= c.cfg.Senders; i++ {
//		go func(instance int) {
//			time.Sleep(time.Duration(instance-1) * time.Duration(c.cfg.LoadInterval) * time.Millisecond)
//			clientSender, err := getClient(ctx, c.cfg.Host, c.cfg.Port, fmt.Sprintf("%s-sender-%d", c.localClientID, instance))
//			if err != nil {
//				c.Logf("error connecting sender: %d, error: %s", instance, err.Error())
//				c.stats.SendErrors.Inc()
//				return
//			}
//			defer func() {
//				err = clientSender.Close()
//				if err != nil {
//					c.Logf("instance:%d, sender: %d, error: %s", c.Id, instance, err.Error())
//					c.stats.SendErrors.Inc()
//				}
//			}()
//			cycles := 0
//			sendCh := make(chan *kubemq.EventStore, c.cfg.SendBatch)
//			receiveCh := make(chan *kubemq.EventStoreResult, c.cfg.SendBatch)
//			errCh := make(chan error, 1)
//			isStreamUp := false
//			for {
//				if c.cfg.SendCycles > 0 {
//					if cycles >= c.cfg.SendCycles {
//						c.Logf("instance:%d, sender: %d, completed", c.Id, instance)
//						return
//					} else {
//						cycles++
//					}
//				}
//				if !isStreamUp {
//					go clientSender.StreamEventsStore(ctx, sendCh, receiveCh, errCh)
//					isStreamUp = true
//				}
//				select {
//				case <-time.After(time.Duration(c.cfg.SendInterval) * time.Second):
//					go func() {
//						//defer func() {
//						//	sendDone <- true
//						//}()
//						randID := nuid.New().Next()
//						for j := 0; j < c.cfg.SendBatch; j++ {
//							msgID := fmt.Sprintf("client-%d-%d-%s", c.Id, instance, randID)
//							event := clientSender.ES().SetId(msgID).SetChannel(c.localChannel).SetBody(randomMessageBody)
//							select {
//							case sendCh <- event:
//							case <-ctx.Done():
//								return
//							default:
//
//							}
//						}
//					}()
//					start := time.Now()
//					cnt := 0
//					for {
//						select {
//						case result := <-receiveCh:
//							cnt++
//							if result.Sent {
//								c.stats.Messages.Inc()
//							} else {
//								c.Log(result.Err.Error())
//								c.stats.SendErrors.Inc()
//							}
//							if cnt >= c.cfg.SendBatch {
//								goto done
//							}
//						case err := <-errCh:
//							c.stats.SendErrors.Inc()
//							c.Logf("instance:%d, sender: %d, error: %s", c.Id, instance, err.Error())
//							isStreamUp = false
//							goto done
//						case <-ctx.Done():
//							goto done
//						}
//
//					}
//				done:
//					end := time.Since(start)
//					c.stats.SendLatency.Store(end.Nanoseconds())
//
//				case <-ctx.Done():
//					return
//				}
//			}
//
//		}(i)
//	}
//	<-ctx.Done()
//}
//
//func (c *StoreClient) GetClientStats() *ClientStats {
//	return c.stats
//}
