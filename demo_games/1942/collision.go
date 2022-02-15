package main

import "github.com/jf-tech/console/cgame"

func registerCollidable(sp *cgame.SpriteManager) {
	sp.CollidableRegistry().
		RegisterBulk(
			alphaName,
			[]string{
				betaName,
				betaBulletName,
				gammaName,
				gammaBulletName,
				deltaName,
				bossName,
				bossBulletName,
				giftPackName}).
		RegisterBulk(
			alphaBulletName,
			[]string{
				betaName,
				gammaName,
				deltaName,
				bossName,
			})
}
