ALTER TABLE orders DROP column discount;
ALTER TABLE orders ADD column shippingId INTEGER references shippingOptions(id);
ALTER TABLE orders ADD column discountId INTEGER references discounts(id);