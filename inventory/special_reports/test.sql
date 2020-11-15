	Select 
				item_name,
				STUFF((
SELECT ' ' + specs
FROM #temp
where item_name = table_temp.item_name
FOR XML PATH (''))
,1,1,'') as specs
,format(sum(convert(int,qty)),'##,###') qty
,format((convert(money,uprice)),'###,####,###.00') uprice
,format(sum(convert(money,tprice)),'###,###,###.00') tprice
				from 
				#temp table_temp
				group by item_name,uprice
				order by 1
	end 