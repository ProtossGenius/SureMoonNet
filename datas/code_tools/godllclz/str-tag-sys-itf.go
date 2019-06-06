package pgt_interface

type StrTagSysItf interface {
	AddTag(tagval, node string) (tagid int)
	/*
	 * tag status may drop/use etc.
	 * status is a tag_id.
	 * result means suc or fail
	 */
	PutTagStatus(tagid, status int) (result bool)
	PutTwoTagRelation(taga, tagb, relation int)
	/*
	 * same to PutTwoTagRelation, but it is not sure put. some days may drop this function.
	 */
	MayTowTagRealtion(taga, tagb, relation int)
	/*
	 * if relation < 0, means taga, tagb not found, but tagb, taga found.
	 */
	GetTwoTagRelation(taga, tagb int) (relation int)
	GetTagsetRelation(tagid int) (relationJson string)
	PutTagsetRelation(tagid int, relationJson string)
	//count means return array's len will not more than count
	GetFathers(count, tagid int) (fatherid []int)
	GetSons(count, tagid int) (sons []int)
	GetTagsets(num, id int) (sets []int)
	//get tagsets only from fsets.
	GetTagsetsFrom(Num, tagid int, fsets []int) (gsets []int)
}
