package bootstrap

func (b *Bootstrap) BuildRestBootstrap() *Bootstrap {
	// set dependencies
	b.setConfig()
	b.setPostgresDB()
	b.setRestApiRepository()
	b.setLogger()
	b.setJWTToken()
	b.setCache()
	b.setRabbitMQ()

	return b
}

func (b *Bootstrap) BuildConsumerUpdateOrderStatusBootstrap() *Bootstrap {
	// set dependencies
	b.setConfig()
	b.setCache()
	b.setPostgresDB()
	b.SetUpdateStatusConsumerRepository()
	b.setLogger()
	b.setRabbitMQ()

	return b
}

func (b *Bootstrap) BuildConsumerExpiredPaymentBootstrap() *Bootstrap {
	// set dependencies
	b.setConfig()
	b.setCache()
	b.setPostgresDB()
	b.SetExpiredPaymentConsumerRepository()
	b.setLogger()
	b.setRabbitMQ()

	return b
}
