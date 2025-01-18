drop view owner_with_books ;
create view owner_with_books as 
select
	o.*,
	(
	select
		array_to_json( array_agg(book_list.*))as books
	from
		(
		select
			b.*
		from
			books b
		where
			b.owner_id = o.id
) as book_list
)as Books
from
	"lender" o ;
