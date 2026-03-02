-- MySQL dump 10.13  Distrib 8.0.41, for Win64 (x86_64)
--
-- Host: localhost    Database: mathovercomedb
-- ------------------------------------------------------
-- Server version	8.0.41

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `categories`
--

DROP TABLE IF EXISTS `categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `categories` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `categories`
--

LOCK TABLES `categories` WRITE;
/*!40000 ALTER TABLE `categories` DISABLE KEYS */;
INSERT INTO `categories` VALUES (1,'数と式'),(2,'2次関数'),(3,'図形と計量'),(4,'データの分析'),(5,'確率'),(6,'図形の性質'),(7,'整数');
/*!40000 ALTER TABLE `categories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `choices`
--

DROP TABLE IF EXISTS `choices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `choices` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `problem_id` bigint NOT NULL,
  `choice_text` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `is_correct` tinyint(1) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `problem_id` (`problem_id`),
  CONSTRAINT `fk_choices_problem` FOREIGN KEY (`problem_id`) REFERENCES `problems` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=655 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `choices`
--

LOCK TABLES `choices` WRITE;
/*!40000 ALTER TABLE `choices` DISABLE KEYS */;
INSERT INTO `choices` VALUES (487,197,'x⁴ - 5x³ + 18x² - 15x + 25',0),(488,197,'x⁴ - 10x³ + 35x² - 50x + 24',1),(489,197,'x⁴ + 16x³ + 8x² - 50x + 20',0),(490,197,'x⁴ + 24x² - 24x² - 36x + 12',0),(491,198,'(x-1)(x+5)(x+3)(x-2)',0),(492,198,'(x+1)(x+2)(x-1)(x+4)',1),(493,198,'(x-1)(x-6)(x+2)(x+3)',0),(494,198,'(x+1)(x-2)(x+3)(x-7)',0),(495,199,'(a-b)(b+c)(c-a)',0),(496,199,'-(a-b)(b-c)(c-a)',1),(497,199,'(a-b)(b-c)(c-a)',0),(498,199,'-(a-b)(b-c)(c+a)',0),(499,200,'(2x²-1x+1)(2x²-2x+1)',0),(500,200,'(2x²-2x+1)(2x²-2x+1)',0),(501,200,'(2x²+2x+1)(2x²-2x+1)',1),(502,200,'(2x²+2x+1)(2x²+2x+1)',0),(503,201,'x=1,4',0),(504,201,'x=1,3',1),(505,201,'x=-1,2',0),(506,201,'x=5,8',0),(507,202,'x<5',0),(508,202,'x>1',1),(509,202,'x>8',0),(510,202,'x<3',0),(511,203,'k=3',0),(512,203,'k=-4',1),(513,203,'k=-1',0),(514,203,'k=2',0),(515,204,'a=2,b=3 または a=-2,b=7',1),(516,204,'a=2,b=-3 または a=-2,b=7',0),(517,204,'a=2,b=3 または a=2,b=-7',0),(518,204,'a=2,b=3 または a=-2,b=-7',0),(519,205,'5',0),(520,205,'2',1),(521,205,'4',0),(522,205,'1',0),(523,206,'2<a<5',0),(524,206,'4<a<8',0),(525,206,'a<-4 または a>5',1),(526,206,'a<-2 または a>3',0),(527,207,'-2<m<3',0),(528,207,'-1<m<2',0),(529,207,'3<m<4',1),(530,207,'6<m<9',0),(531,208,'a=-4,b=4',0),(532,208,'a=-2,b=3',0),(533,208,'a=-1,b=2',1),(534,208,'a=-3,b=1',0),(535,209,'60°,120°,180°',0),(536,209,'45°,90°,135°',0),(537,209,'30°,90°,150°',1),(538,209,'30°,60°,90°',0),(539,210,'0°<θ<30°',0),(540,210,'30°<θ<150°',1),(541,210,'45°<θ<135°',0),(542,210,'0°<θ<45°',0),(543,211,'27',0),(544,211,'24',0),(545,211,'36',1),(546,211,'45',0),(547,212,'(2+√3)/2',0),(548,212,'(2+√3)/3',0),(549,212,'(4-√7)/3',1),(550,212,'(4-√7)/2',0),(551,213,'5√3/2',0),(552,213,'2√2/5',0),(553,213,'5√2/8',1),(554,213,'3√2/8',0),(555,214,'0°≤θ≤30°、θ=180°',0),(556,214,'0°≤θ≤60°、θ=180°',1),(557,214,'0°≤θ≤120°、θ=180°',0),(558,214,'0°≤θ≤135°、θ=180°',0),(559,215,'4',0),(560,215,'3',0),(561,215,'8',0),(562,215,'6',1),(563,216,'12',0),(564,216,'3',0),(565,216,'1',1),(566,216,'2',0),(567,217,'0.68',0),(568,217,'0.88',1),(569,217,'0.58',0),(570,217,'0.45',0),(571,218,'増加',0),(572,218,'この条件では分からない',0),(573,218,'減少',1),(574,218,'一致',0),(575,219,'36',0),(576,219,'34',0),(577,219,'37',1),(578,219,'41',0),(579,220,'0.81',0),(580,220,'0.77',1),(581,220,'0.68',0),(582,220,'0.65',0),(583,221,'140',0),(584,221,'135',1),(585,221,'125',0),(586,221,'130',0),(587,222,'31',0),(588,222,'28',0),(589,222,'44',1),(590,222,'24',0),(591,223,'65',0),(592,223,'50',0),(593,223,'55',1),(594,223,'70',0),(595,224,'32/55',0),(596,224,'27/55',1),(597,224,'18/55',0),(598,224,'6/55',0),(599,225,'1/5',0),(600,225,'3/10',0),(601,225,'2/5',1),(602,225,'1/6',0),(603,226,'5/13',0),(604,226,'3/13',0),(605,226,'6/13',1),(606,226,'4/13',0),(607,227,'18/7',0),(608,227,'27/7',0),(609,227,'21/4',1),(610,227,'24/5',0),(611,228,'2',0),(612,228,'6',0),(613,228,'3',1),(614,228,'4',0),(615,229,'20:3',0),(616,229,'18:7',0),(617,229,'25:7',1),(618,229,'16:3',0),(619,230,'2/9',0),(620,230,'7/9',1),(621,230,'2/3',0),(622,230,'7/3',0),(623,231,'8:3',0),(624,231,'4:3',1),(625,231,'4:1',0),(626,231,'6:1',0),(627,232,'2/7',0),(628,232,'1/7',1),(629,232,'2/9',0),(630,232,'1/5',0),(631,233,'60',0),(632,233,'80',0),(633,233,'70',1),(634,233,'50',0),(635,234,'n=2,6',0),(636,234,'n=4,1',0),(637,234,'n=7,1',1),(638,234,'n=5,7',0),(639,235,'8',0),(640,235,'12',0),(641,235,'6',1),(642,235,'4',0),(643,236,'4',0),(644,236,'6',1),(645,236,'1',0),(646,236,'5',0),(647,237,'57',0),(648,237,'53',0),(649,237,'67',1),(650,237,'76',0),(651,238,'32',0),(652,238,'53',1),(653,238,'47',0),(654,238,'28',0);
/*!40000 ALTER TABLE `choices` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `problems`
--

DROP TABLE IF EXISTS `problems`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `problems` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `category_id` int NOT NULL,
  `question` longtext,
  `hint` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `answer` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `category_id` (`category_id`),
  CONSTRAINT `problems_ibfk_1` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=239 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `problems`
--

LOCK TABLES `problems` WRITE;
/*!40000 ALTER TABLE `problems` DISABLE KEYS */;
INSERT INTO `problems` VALUES (197,1,'次の式を展開せよ<br>(x-1)(x-2)(x-3)(x-4)','計算する順番を工夫しましょう！',NULL),(198,1,'次の式を因数分解せよ<br>(x²+3x)² - 2(x²+3x) - 8','x²+3xをひとまとめにして考えてみましょう (t = x² + 3x とおくなど)',NULL),(199,1,'次の式を因数分解せよ<br>(b−c)a² + (c−a)b² + (a−b)c²','次数の低い文字（どれでもよい）に着目して整理してみましょう！',NULL),(200,1,'次の式を因数分解せよ<br>4x⁴ + 1','a⁴ + b⁴ の因数分解は (a² + b²)² − 2a²b² の展開を考えましょう！',NULL),(201,1,'次の方程式を解け<br>|x−1| + |x−2| = x','絶対値の中が正か負かで場合分けしましょう！',NULL),(202,1,'次の不等式を解け<br>|x−4| < 3x','絶対値の中が正か負かで場合分けしましょう！',NULL),(203,2,'y = −2x² + 8x + k の最大値が 4 であるとき k の値を求めよ','関数の頂点の y 座標を求めてみましょう！',NULL),(204,2,'0 ≤ x ≤ 3 のとき f(x)=ax²−2ax+b の最大値が 9、最小値が 1 のとき a,b を求めよ','a=0, a<0, a>0 の場合で場合分けしてみましょう！',NULL),(205,2,'x + 2y = 3 のとき 2x² + y² の最小値を求めよ','x=3−2y を代入して y だけの式で考えましょう！',NULL),(206,2,'2次方程式 ax²−(a+1)x−(a+3)=0 が −1 < x < 0, 1 < x < 2 でそれぞれ 1 つの実数解をもつとき、a の範囲を求めよ','x=−1,0,1,2 での式の符号変化を調べましょう！',NULL),(207,2,'y = x²−mx+m²−3m のグラフが x 軸の正の部分と異なる 2 点で交わるときの m の範囲を求めよ','判別式 D>0, 軸の x 座標>0, f(0)>0 の 3 つを満たすように条件を立てましょう！',NULL),(208,2,'不等式 ax² + bx + 3 > 0 の解が −1 < x < 3 であるとき a,b を求めよ','放物線が上向きか下向きか、また x=−1,3 での値を調べましょう！',NULL),(209,3,'2cos²θ + 3sinθ − 3 = 0 (0°≤θ≤180°) を解け','cos²θ = 1−sin²θ に直して sinθ の 2 次方程式にしましょう！',NULL),(210,3,'0°≤θ≤180° のとき sinθ > 1/2 を満たす θ の範囲を求めよ','単位円上で y = 1/2 より上の角度を考えましょう！',NULL),(211,3,'円に内接する四角形 ABCD がある。AB=4, BC=5, CD=7, DA=10 のとき cos A の値を求め、それを利用して四角形の面積を求めよ','△ABC と △BCD に分けて面積を求め、sin(180°−A)=sin A を使いましょう！',NULL),(212,3,'cosθ − sinθ = 1/2 (0°<θ<180°) のとき tan θの値を求めよ','sinθ, cosθ を求めて tanθ=sinθ/cosθ で計算しましょう！',NULL),(213,3,'sinθ + cosθ = √2/2 (0°<θ<180°) のとき sin³θ + cos³θ の値を求めよ','a³ + b³ = (a + b)(a² − ab + b²) を使いましょう！',NULL),(214,3,'0°≤θ≤180° のとき 2sin²θ − cosθ − 1 ≤ 0 の不等式を解け','sin²θ = 1−cos²θ に代えて cosθ の不等式にしましょう！',NULL),(215,4,'次のデータ {5,7,4,3,6} における分散を求めよ','分散 = x² の平均 − (x の平均)² の公式を使いましょう！',NULL),(216,4,'次のデータ {5,4,8,12,17,24,27,28,22,30,9,6} で 30 が 18 だったとき、平均は修正前よりいくつ減少するか','平均値 = 総和 ÷ データ数 で計算しましょう！',NULL),(217,4,'50点満点のテスト A,B を行った結果の得点<br><br><table border=\"1\"><tr><th>生徒</th><th>1</th><th>2</th><th>3</th><th>4</th><th>5</th><th>6</th><th>7</th><th>8</th><th>9</th><th>10</th></tr><tr><td>x</td><td>43</td><td>41</td><td>43</td><td>38</td><td>39</td><td>42</td><td>42</td><td>39</td><td>41</td><td>42</td></tr><tr><td>y</td><td>49</td><td>42</td><td>44</td><td>36</td><td>40</td><td>44</td><td>45</td><td>42</td><td>42</td><td>46</td></tr></table><br>このとき相関係数 r を小数第3位で四捨五入して求めよ','相関係数の公式を用いて計算しましょう！',NULL),(218,4,'次のデータ {5,4,8,12,17,24,27,28,22,30,9,6} で 6→10, 30→26 に修正したとき、分散は修正前よりどうなるか','平均値を求め、修正した 2 つの数だけで偏差平方和を比べましょう！',NULL),(219,4,'あるクラスのテスト平均が 54.3 点で、得点が 69,65,62,57,55,55,53,48,42,x のとき x の値を求めよ','x を用いて (既存の合計 + x) ÷ 10 = 54.3 の式を立てましょう！',NULL),(220,4,'30点満点のテスト A,B を行った結果の得点<br><br><table border=\"1\"><tr><th>生徒</th><th>1</th><th>2</th><th>3</th><th>4</th><th>5</th><th>6</th><th>7</th><th>8</th><th>9</th><th>10</th></tr><tr><td>x</td><td>29</td><td>25</td><td>22</td><td>28</td><td>18</td><td>23</td><td>26</td><td>30</td><td>30</td><td>29</td></tr><tr><td>y</td><td>23</td><td>23</td><td>18</td><td>26</td><td>17</td><td>20</td><td>21</td><td>20</td><td>26</td><td>26</td></tr></table><br>このとき相関係数 r を小数第3位で四捨五入して求めよ','相関係数の公式を用いて計算しましょう！',NULL),(221,5,'大、中、小 3 個のサイコロを投げるとき、目の積が 4 の倍数になる場合は何通りか','(全体) – (積が4の倍数でない場合) で考えましょう！',NULL),(222,5,'5 人に招待状を送るため、宛名を書いた招待状と封筒を作成した。招待状を全部無作為に封筒に入れたとき、誰も自分の封筒に入らない場合は何通りか','1～5 の並びで、k 番目が k にならない並べ方を数えましょう！',NULL),(223,5,'x + y + z = 9, x≥0, y≥0, z≥0 を満たす整数の組は何通りか','9 個の○と 2 本の仕切りを置く方法を考えましょう！',NULL),(224,5,'赤、青、黄の札がそれぞれ 4 枚ずつあり、各札に 1～4 の番号が書かれている。12 枚の札から 3 枚取り出すとき、番号がすべて異なる確率を求めよ','番号を選ぶ方法 × 色を選ぶ方法 で考えましょう！',NULL),(225,5,'袋の中に赤球2個、白球3個がある。A,B が交互に1個ずつ取り出し、2 個目の赤球を取った方が勝ちとする。取り出した球は戻さない。B が勝つ確率を求めよ','B が勝つパターンを全部書き、確率を足しましょう！',NULL),(226,5,'工場では製品を機械A(不良4%)とその他(不良7%)で作る。全体の60%をA製とするとき、不良品だったものがA製である確率を求めよ','条件付き確率：(Aの不良率×Aの割合)÷全体の不良率 で求めましょう！',NULL),(227,6,'AB=7, BC=5, CA=3 の△ABCにおいて、角Aおよびその外角の二等分線が辺BCまたはその延長と交わる点をそれぞれD,Eとする。線分DEの長さを求めよ','内角・外角の二等分線の定理を使いましょう！',NULL),(228,6,'右図の△ABCで、D,Eはそれぞれ辺BC,CAの中点。ADとBEの交点をF、AFの中点をG、CGとBEの交点をHとする。BE=9のとき△EBCと△FBDの面積比を求めよ','高さが共通なので底辺の長さ比で面積比を求めましょう！',NULL),(229,6,'△ABCの辺BC,CA,ABを3:2に内分する点をそれぞれD,E,Fとする。△ABCと△DEFの面積の比を求めよ','内分比から相似比を考えましょう！',NULL),(230,6,'1辺の長さが7の正三角形ABCがある。AB上にAD=3、AC上にAE=6となるようにD,Eをとる。このときBE,CDの交点をF、直線AFとBCの交点をGとするとき線分CGの長さを求めよ','チェバの定理を使いましょう！',NULL),(231,6,'△ABCにおいて、辺AB上と辺ACの延長上にE,FをとりAE:EB=1:2, AF:FC=3:1とする。直線EFと直線BCの交点をDとするときBD:DCを求めよ','メネラウスの定理を利用しましょう！',NULL),(232,6,'面積が1の△ABCにおいて、辺BC,CA,ABを2:1に内分する点をそれぞれL,M,Nとし、ALとBM, BMとCN, CNとALの交点をそれぞれP,Q,Rとするとき△PQRの面積を求めよ','メネラウスの定理でAP:PR:RLを求めましょう！',NULL),(233,7,'√(63n/40) が有理数となる最小の自然数 n を求めよ','素因数分解して √ が外れる条件を考えましょう！',NULL),(234,7,'√(n²+15) が自然数となる自然数 n をすべて求めよ','√(n²+15)=m とおいて、両辺を2乗して考えましょう！',NULL),(235,7,'25! を計算すると末尾には 0 が何個連続して並ぶか求めよ','25 ÷ 5 ＋ 25 ÷ 25 で 5 の個数を数えましょう！',NULL),(236,7,'整数 a を 7 で割ると 3 余るとき、a¹⁰⁰⁰ を 7 で割ったときの余りを求めよ','3 を何回か掛けていくと、余りに繰り返しのパターンがあることに気づきます。その性質を利用しましょう！',NULL),(237,7,'7n+4 と 8n+5 が互いに素になるような 100 以下の自然数 n は何個あるか','8n+5 − (7n+4) = n+1 に注目して、その数が 1 以外の約数を持たない n を探してみましょう！',NULL),(238,7,'3 で割ると 2 余り、5 で割ると 3 余り、7 で割ると 4 余る自然数 n の最小値を求めよ','「3 で割ると 2 余る数」「5 で割ると 3 余る数」などを順に調べて、すべての条件を満たす数を探してみましょう！',NULL);
/*!40000 ALTER TABLE `problems` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sessionproblems`
--

DROP TABLE IF EXISTS `sessionproblems`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sessionproblems` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `session_id` bigint NOT NULL,
  `problem_id` bigint NOT NULL,
  `selected_choice_id` bigint DEFAULT NULL,
  `is_correct` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `session_id` (`session_id`),
  KEY `problem_id` (`problem_id`),
  KEY `selected_choice_id` (`selected_choice_id`),
  CONSTRAINT `fk_sessionproblems_choice` FOREIGN KEY (`selected_choice_id`) REFERENCES `choices` (`id`),
  CONSTRAINT `fk_sessionproblems_problem` FOREIGN KEY (`problem_id`) REFERENCES `problems` (`id`),
  CONSTRAINT `fk_sessionproblems_session` FOREIGN KEY (`session_id`) REFERENCES `testsessions` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1849 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sessionproblems`
--

LOCK TABLES `sessionproblems` WRITE;
/*!40000 ALTER TABLE `sessionproblems` DISABLE KEYS */;
INSERT INTO `sessionproblems` VALUES (1715,143,197,490,0),(1716,143,202,NULL,NULL),(1717,143,204,NULL,NULL),(1718,143,207,NULL,NULL),(1719,143,211,NULL,NULL),(1720,143,214,NULL,NULL),(1721,143,220,NULL,NULL),(1722,143,218,NULL,NULL),(1723,143,222,NULL,NULL),(1724,143,226,NULL,NULL),(1725,143,227,NULL,NULL),(1726,143,232,NULL,NULL),(1727,144,202,510,0),(1728,144,197,NULL,NULL),(1729,144,208,NULL,NULL),(1730,144,206,NULL,NULL),(1731,144,210,NULL,NULL),(1732,144,213,NULL,NULL),(1733,144,215,NULL,NULL),(1734,144,217,NULL,NULL),(1735,144,225,NULL,NULL),(1736,144,226,NULL,NULL),(1737,144,231,NULL,NULL),(1738,144,228,NULL,NULL),(1739,145,201,NULL,NULL),(1740,145,198,NULL,NULL),(1741,145,203,NULL,NULL),(1742,145,205,NULL,NULL),(1743,145,214,NULL,NULL),(1744,145,213,NULL,NULL),(1745,145,215,NULL,NULL),(1746,145,216,NULL,NULL),(1747,145,223,NULL,NULL),(1748,145,226,NULL,NULL),(1749,145,229,NULL,NULL),(1750,145,227,NULL,NULL),(1751,146,201,NULL,NULL),(1752,146,199,NULL,NULL),(1753,146,205,NULL,NULL),(1754,146,204,NULL,NULL),(1755,146,210,NULL,NULL),(1756,146,211,NULL,NULL),(1757,146,219,NULL,NULL),(1758,146,216,NULL,NULL),(1759,146,226,NULL,NULL),(1760,146,225,NULL,NULL),(1761,146,230,NULL,NULL),(1762,146,227,NULL,NULL),(1763,147,201,505,0),(1764,147,200,502,0),(1765,147,208,534,0),(1766,147,203,NULL,NULL),(1767,147,213,NULL,NULL),(1768,147,210,542,0),(1769,147,218,NULL,NULL),(1770,147,220,NULL,NULL),(1771,147,224,NULL,NULL),(1772,147,222,NULL,NULL),(1773,147,232,NULL,NULL),(1774,147,231,NULL,NULL),(1775,148,201,NULL,NULL),(1776,148,202,NULL,NULL),(1777,148,208,NULL,NULL),(1778,148,207,NULL,NULL),(1779,148,213,NULL,NULL),(1780,148,209,NULL,NULL),(1781,148,220,NULL,NULL),(1782,148,219,NULL,NULL),(1783,148,224,NULL,NULL),(1784,148,226,NULL,NULL),(1785,148,227,NULL,NULL),(1786,148,232,NULL,NULL),(1787,149,201,504,1),(1788,149,199,NULL,NULL),(1789,149,207,530,0),(1790,149,205,NULL,NULL),(1791,149,210,NULL,NULL),(1792,149,214,NULL,NULL),(1793,149,217,NULL,NULL),(1794,149,219,NULL,NULL),(1795,149,223,NULL,NULL),(1796,149,225,NULL,NULL),(1797,149,230,NULL,NULL),(1798,149,229,NULL,NULL),(1799,150,198,494,0),(1800,150,202,509,0),(1801,150,205,520,1),(1802,150,208,533,1),(1803,150,213,554,0),(1804,150,209,538,0),(1805,150,215,NULL,NULL),(1806,150,219,578,0),(1807,150,222,NULL,NULL),(1808,150,223,NULL,NULL),(1809,150,229,NULL,NULL),(1810,150,228,NULL,NULL),(1811,151,199,497,0),(1812,151,198,494,0),(1813,151,208,533,1),(1814,151,204,516,0),(1815,151,214,558,0),(1816,151,213,553,1),(1817,151,219,576,0),(1818,151,217,568,1),(1819,151,223,594,0),(1820,151,224,597,0),(1821,151,232,628,1),(1822,151,230,619,0),(1823,151,234,638,0),(1824,151,236,645,0),(1825,152,198,494,0),(1826,152,202,509,0),(1827,152,204,518,0),(1828,152,208,533,1),(1829,152,211,546,0),(1830,152,214,557,0),(1831,152,216,NULL,NULL),(1832,152,219,578,0),(1833,152,225,602,0),(1834,152,226,604,0),(1835,152,230,622,0),(1836,152,232,630,0),(1837,153,199,498,0),(1838,153,200,501,1),(1839,153,205,522,0),(1840,153,207,528,0),(1841,153,212,547,0),(1842,153,211,546,0),(1843,153,216,565,1),(1844,153,219,575,0),(1845,153,224,598,0),(1846,153,222,588,0),(1847,153,232,629,0),(1848,153,229,615,0);
/*!40000 ALTER TABLE `sessionproblems` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `testsessions`
--

DROP TABLE IF EXISTS `testsessions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `testsessions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint DEFAULT NULL,
  `include_integers` tinyint(1) NOT NULL,
  `start_time` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `fk_testsessions_users` (`user_id`),
  CONSTRAINT `fk_testsessions_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=154 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `testsessions`
--

LOCK TABLES `testsessions` WRITE;
/*!40000 ALTER TABLE `testsessions` DISABLE KEYS */;
INSERT INTO `testsessions` VALUES (143,2,0,'2025-05-15 11:39:21'),(144,2,0,'2025-05-15 11:44:02'),(145,2,0,'2025-05-15 13:49:38'),(146,2,0,'2025-05-15 13:49:46'),(147,2,0,'2025-05-15 13:58:01'),(148,2,0,'2025-05-15 15:17:01'),(149,2,0,'2025-05-15 20:48:38'),(150,2,0,'2025-05-17 06:46:58'),(151,2,1,'2025-05-22 06:44:10'),(152,2,0,'2025-05-22 11:16:03'),(153,2,0,'2025-06-09 19:01:13');
/*!40000 ALTER TABLE `testsessions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_name` varchar(255) DEFAULT NULL,
  `user_id` varchar(255) DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (2,'恭平','kaizokuoukyou','Peace0130'),(3,'dg32523','234543232','Peace0130'),(4,'342343243423','34234434','11111111'),(5,'43242342','32432434324324','11111111'),(6,'435','retrter','aaaaaaaaaa');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-03-01  6:55:53
